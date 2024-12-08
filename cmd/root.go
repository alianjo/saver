package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/alianjo/saver/client"
	"github.com/itaysk/kubectl-neat/cmd"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getLineNumber retrieves the line number of the caller
func getLineNumber() int {
	_, _, line, _ := runtime.Caller(1) // The third return value is the line number
	return line
}

var opt struct {
	namespace string
	output    string
}

// NewBackupCMD creates a new Cobra command for backup
func NewBackupCMD() *cobra.Command {
	rootcmd := &cobra.Command{
		Use:     "save",
		Short:   "sv",
		Long:    "save workloads",
		Example: "kubectl save deployment --namespace controller -o controller-deployments.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// Correctly report the file and line number
				return fmt.Errorf("Pass one of the options [deployment, statefulset, daemonset] is required at %s:%d", "cmd/backup.go", getLineNumber())
			}
			// Use the package-level opt struct for namespace and workload
			return backup(args[0], opt.namespace, opt.output)
		},
		Version: "0.1",
	}

	// Correctly bind the flags to the opt struct
	rootcmd.Flags().StringVarP(&opt.namespace, "namespace", "n", "default", " namespace")
	rootcmd.Flags().StringVarP(&opt.output, "output", "o", "", "output path")

	return rootcmd
}

// backup simulates the backup process for a Kubernetes secret

func de_cluttering(in string, apiVersion string, kind string) (string, error) {
	// Parse the input JSON
	var resource map[string]interface{}
	if err := json.Unmarshal([]byte(in), &resource); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Add apiVersion and kind if they don't exist
	resource["apiVersion"] = apiVersion
	resource["kind"] = kind

	// Convert back to JSON for further processing
	updatedJSON, err := json.Marshal(resource)
	if err != nil {
		return "", fmt.Errorf("error marshalling updated JSON: %v", err)
	}

	// Use the cmd.Neat function to clean up the resource
	out, err := cmd.Neat(string(updatedJSON))
	if err != nil {
		return "", fmt.Errorf("error cleaning resource: %v", err)
	}

	return out, nil
}

func convertToJson(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling data to JSON: %v", err)
	}
	return string(jsonData), nil
}

func PrintWorkloadYaml(clientset *kubernetes.Clientset, namespace string, workload string) {
	ctx := context.Background()

	// Define apiVersion and kind based on the workload type
	var apiVersion, kind string
	switch workload {
	case "deployment":
		apiVersion = "apps/v1"
		kind = "Deployment"
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching deployments: %v\n", err)
			return
		}
		for _, d := range deployments.Items {
			json, err := convertToJson(d)
			if err != nil {
				fmt.Printf("Error converting deployments to JSON: %v\n", err)
				return
			}
			cleanData, err := de_cluttering(string(json), apiVersion, kind)
			if err != nil {
				fmt.Printf("Error cleaning JSON: %v\n", err)
				return
			}
			printYaml(cleanData)
		}
	case "daemonset":
		apiVersion = "apps/v1"
		kind = "DaemonSet"
		daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching daemonsets: %v\n", err)
			return
		}
		for _, d := range daemonSets.Items {
			json, err := convertToJson(d)
			if err != nil {
				fmt.Printf("Error converting daemonSets to JSON: %v\n", err)
				return
			}
			cleanData, err := de_cluttering(string(json), apiVersion, kind)
			if err != nil {
				fmt.Printf("Error cleaning JSON: %v\n", err)
				return
			}
			printYaml(cleanData)
		}
	case "statefulset":
		apiVersion = "apps/v1"
		kind = "StatefulSet"
		statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error fetching statefulsets: %v\n", err)
			return
		}
		for _, d := range statefulSets.Items {
			json, err := convertToJson(d)
			if err != nil {
				fmt.Printf("Error converting statefulSets to JSON: %v\n", err)
				return
			}
			cleanData, err := de_cluttering(string(json), apiVersion, kind)
			if err != nil {
				fmt.Printf("Error cleaning JSON: %v\n", err)
				return
			}
			printYaml(cleanData)
		}
	default:
		fmt.Printf("Invalid workload type: %s\n", workload)
	}
}

// PrintWorkloadYamlToFile fetches the YAML of the specified workload resources and writes them to a file.
func PrintWorkloadYamlToFile(clientset *kubernetes.Clientset, namespace string, workload string, filename string) error {
	ctx := context.Background()
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Define apiVersion and kind based on the workload type
	var apiVersion, kind string
	switch workload {
	case "deployment":
		apiVersion = "apps/v1"
		kind = "Deployment"
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching deployments: %v", err)
		}
		for _, d := range deployments.Items {
			json, err := convertToJson(d)
			if err != nil {
				return fmt.Errorf("error converting deployments to JSON: %v", err)
			}
			cleanData, err := de_cluttering(string(json), apiVersion, kind)
			if err != nil {
				return fmt.Errorf("error cleaning JSON: %v", err)
			}
			if err := writeYamlToFile(file, cleanData); err != nil {
				return err
			}
		}
	// Similar blocks for daemonset and statefulset...
	default:
		return fmt.Errorf("invalid workload type: %s", workload)
	}

	return nil
}

func writeYamlToFile(file *os.File, resource string) error {
	var data interface{}
	err := yaml.Unmarshal([]byte(resource), &data)
	if err != nil {
		return fmt.Errorf("Error unmarshalling resource to YAML: %v\n", err)
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling resource to YAML: %v", err)
	}
	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}
	_, err = file.WriteString("\n---\n") // Separate resources in the file with ---
	return err
}

// A function to Add
// func dataAdder() {}

func printYaml(resource string) {
	var data interface{}
	err := yaml.Unmarshal([]byte(resource), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling resource to YAML: %v\n", err)
		return
	}
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshalling resource to YAML: %v\n", err)
		return
	}
	fmt.Println("---")
	fmt.Println(string(yamlData))
}

func backup(workload, namespace, output string) error {
	client := client.K8Client()

	if output != "" {
		PrintWorkloadYamlToFile(client, namespace, workload, output)
	} else {
		PrintWorkloadYaml(client, namespace, workload)
	}

	return nil
}
