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

	options := k8s.NewKubectlOptions("", "", "ocdp") //put your name spaces
	decision := k8s.AreAllNodesReady(t, options)
	color.Set(color.FgGreen)
	fmt.Println("Is Kubernetes up and running ?", decision)
	color.Unset()
	k8s.RunKubectlAndGetOutputE(t, options, "cluster-info")
	k8s.RunKubectlAndGetOutputE(t, options, "get", "nodes")
	//k8s.RunKubectlE(t, options, "get", "pods", "-o=wide", "--field-selector", "status.phase=Running")
	//Checking Pods here
	podsare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].metadata.name}")
	podsstatus, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].status.phase}")
	podsstatusb, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods", "-o=jsonpath={.items[*].status.containerStatuses[0].ready}")
	arraypodb := strings.Split(podsstatusb, " ")
	statuscount := strings.Count(podsstatus, " ")
	kv := strconv.FormatInt(int64(statuscount), 10)
	jv, _ := strconv.Atoi(kv)
	statuscountnum := jv + 1
	//statuscountrunning := strings.Count(podsstatus, "Running")
	arraystatus := strings.Split(podsstatus, " ")
	color.Set(color.FgYellow)
	fmt.Println("checking pod status", podsstatus)
	color.Unset()
	color.Set(color.FgGreen)
	arraypod := strings.Split(podsare, " ")
	runningpods := []string{}
	notrunningpods := []string{}
	//fmt.Println("number of pods running are ", statuscountrunning)
	counter := 0
	color.Unset()
	for j := 0; j <= (len(arraypod) - 1); j++ {
		if arraystatus[j] == "Running" {
			runningpods = append(runningpods, arraypod[j])
			counter = counter + 1
		}
		if arraystatus[j] != "Running" {
			notrunningpods = append(notrunningpods, arraypod[j])
		}
	}
	actualrunningpods := []string{""}
	r := 0
	for k := 0; k <= (counter - 1); k++ {
		if arraypodb[k] == ("true") {
			fmt.Println("Actual running pods online", runningpods[k])
			r = r + 1
			actualrunningpods = append(actualrunningpods, runningpods[k])
		}
		if arraypodb[k] != ("true") {
			notrunningpods = append(notrunningpods, runningpods[k])
		}
	}
	color.Set(color.FgGreen)
	fmt.Println("total pods count", statuscountnum)
	fmt.Println("Total Pods functioning properly =", r)
	if len(arraypod) == r {
		fmt.Println("All Pods are Up and Running")
	}
	color.Unset()
	color.Set(color.FgBlue)
	//k8s.RunKubectlE(t, options, "get", "pods", "--field-selector", "status.phase=Running")
	//checking services and loadbalancers
	servicesare, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", "-o=jsonpath={.items[*].metadata.name}")
	color.Unset()
	fmt.Println("services are", servicesare)
	arrayservice := strings.Split(servicesare, " ")
	a := len(arrayservice) - 1

	for b := 0; b <= a; b++ {
		service := k8s.GetService(t, options, arrayservice[b])
		externalip, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", arrayservice[b], "-o=jsonpath={.status.loadBalancer.ingress[*].ip}")
		targetport, _ := k8s.RunKubectlAndGetOutputE(t, options, "get", "svc", arrayservice[b], "-o", "jsonpath={.spec.ports[*].port}")
		arrayport := strings.Split(targetport, " ")
		c := len(arrayport) - 1

		for d := 0; d <= c; d++ {
			if arrayservice[b] == "kubernetes" {
				goto LOOP
			}
			if externalip == "" {
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
			if arrayport[d] == ("8080") || arrayport[d] == ("80") || arrayport[d] == ("8008") || arrayport[d] == ("591") {
				tlsConfig := tls.Config{}
				http_helper.HttpGetWithRetryWithCustomValidation(
					t,
					fmt.Sprintf("http://%s", endpoint),
					&tlsConfig,
					10,
					10*time.Second,
					func(statusCode int, body string) bool {
						return statusCode == 200
					},
				)
			}
			color.Set(color.FgGreen)
			bar := progressbar.Default(100) //added loading bar
			for i := 0; i < 100; i++ {
				bar.Add(1)
				time.Sleep(5 * time.Millisecond)
			}
			fmt.Println("The test passed Successfully with the", endpoint)
			color.Unset()
		}
	LOOP:
	}
	//list of pods not working
	color.Set(color.FgRed)
	fmt.Println("\n", "The Following Pods need Attention !!")
	for y := 0; y <= (len(notrunningpods) - 1); y++ {
		fmt.Println(notrunningpods[y])
	}
	color.Unset()

}
