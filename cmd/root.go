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

func de_cluttering(in string) (string, error) {
	out, err := cmd.Neat(in)
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
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

	switch workload {
	case "deployment":
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
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}

			printYaml(clean_data)

		}
	case "daemonset":
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
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}

			printYaml(clean_data)

		}
	case "statefulset":
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
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}

			printYaml(clean_data)

		}
	default:
		fmt.Printf("Invalid workload type: %s\n", workload)
	}
}

// PrintWorkloadYamlToFile fetches the YAML of the specified workload resources and writes them to a file.
func PrintWorkloadYamlToFile(clientset *kubernetes.Clientset, namespace string, workload string, filename string) error {
	ctx := context.Background()
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()
	switch workload {
	case "deployment":
		deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching deployments: %v", err)
		}
		for _, d := range deployments.Items {
			json, err := convertToJson(d)
			if err != nil {
				fmt.Printf("Error converting daemonSets to JSON: %v\n", err)
				return err
			}
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}
			if err := writeYamlToFile(file, clean_data); err != nil {
				return err
			}
		}
	case "daemonset":
		daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching daemonsets: %v", err)
		}
		for _, ds := range daemonSets.Items {
			json, err := convertToJson(ds)
			if err != nil {
				fmt.Printf("Error converting daemonSets to JSON: %v\n", err)
				return err
			}
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}
			if err := writeYamlToFile(file, clean_data); err != nil {
				return err
			}
		}
	case "statefulset":
		statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("error fetching statefulsets: %v", err)
		}
		for _, ss := range statefulSets.Items {
			json, err := convertToJson(ss)
			if err != nil {
				fmt.Printf("Error converting daemonSets to JSON: %v\n", err)
				return err
			}
			clean_data, err := de_cluttering(string(json))
			if err != nil {
				fmt.Errorf("error Cleaning Json to file: %v", err)
			}
			if err := writeYamlToFile(file, clean_data); err != nil {
				return err
			}
		}
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

	PrintWorkloadYaml(client, namespace, workload)
	if output != "" {
		PrintWorkloadYamlToFile(client, namespace, workload, output)
	}

	return nil
}
