package test

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/reiver/go-telnet"
	"github.com/schollz/progressbar/v3"
)

func TestKubernetes(t *testing.T) {

	options := k8s.NewKubectlOptions("", "", "default")
	decision := k8s.AreAllNodesReady(t, options)
	fmt.Println("Is Kubernetes up and running ?", decision)
	k8s.RunKubectlAndGetOutputE(t, options, "cluster-info")
	k8s.RunKubectlAndGetOutputE(t, options, "get", "nodes")
	//k8s.RunKubectlE(t, options, "get", "pods", "-o=wide", "--field-selector", "status.phase=Running")
	podsare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].metadata.name}")
	podsstatus, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].status.phase}")
	statuscount := strings.Count(podsstatus, " ")
	kv := strconv.FormatInt(int64(statuscount), 10)
	jv, _ := strconv.Atoi(kv)
	statuscountnum := jv + 1
	fmt.Println("total pods count", statuscountnum)
	statuscountrunning := strings.Count(podsstatus, "Running")
	arraystatus := strings.Split(podsstatus, " ")
	color.Set(color.FgYellow)
	fmt.Println("checking node status", podsstatus)
	color.Unset()
	color.Set(color.FgGreen)
	arraypod := strings.Split(podsare, " ")
	bar := progressbar.Default(100)  //added loading bar 
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(20 * time.Millisecond)
	}
	if len(arraypod) == statuscountrunning {
		fmt.Println("All Pods are Up and Running")
	}
	notrunningpods := []string{" "}
	fmt.Println("number of pods running are ", statuscountrunning)
	color.Unset()
	for j := 0; j <= (len(arraypod) - 1); j++ {
		if arraystatus[j] != "Running" {
			notrunning := arraypod[j]
			fmt.Println("The node is offline", notrunning)
			notrunningpods = append(notrunningpods, notrunning)
		}
	}
	color.Set(color.FgBlue)
	k8s.RunKubectlE(t, options, "get", "pods", "--field-selector", "status.phase=Running")
	servicesare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-o=jsonpath={.items[*].metadata.name}")
	color.Unset()
	fmt.Println("services are", servicesare)
	arrayservice := strings.Split(servicesare, " ")
	a := len(arrayservice) - 1

	for b := 0; b <= a; b++ {
		service := k8s.GetService(t, options, arrayservice[b])
		//externalip, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", arrayservice[b], "-o=jsonpath={.status.loadbalancer.ingress[*].hostname}")
		targetport, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", arrayservice[b], "-o", "jsonpath={.spec.ports[*].targetPort}")

		arrayport := strings.Split(targetport, " ")
		c := len(arrayport) - 1
		for d := 0; d <= c; d++ {
			if arrayservice[b] == "kubernetes" {
				goto LOOP
			}
			e, _ := strconv.Atoi(arrayport[d])
			endpoint := k8s.GetServiceEndpoint(t, options, service, e)
			color.Set(color.FgGreen)
			fmt.Println(endpoint)
			fmt.Println("test connection with ", endpoint)
			color.Unset()
			var caller telnet.Caller = telnet.StandardCaller
			telnet.DialToAndCall(endpoint, caller)
			tlsConfig := tls.Config{}
			http_helper.HttpGetWithRetryWithCustomValidation(
				t,
				fmt.Sprintf("http://%s", endpoint),
				&tlsConfig,
				30,
				10*time.Second,
				func(statusCode int, body string) bool {
					return statusCode == 200
				},
			)
		}
	LOOP:
	}
	if len(notrunningpods) > 1 {
		color.Set(color.FgRed)
		fmt.Println("the following pods are offline \n", notrunningpods)
		color.Unset()
	}
}
