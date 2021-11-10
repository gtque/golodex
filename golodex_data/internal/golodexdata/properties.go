package golodexdata

import (
	"context"
	"log"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

func Property(key string, defaultValue string)string {
	return FullProperty(key, defaultValue, false)
}

func Secret(key string)string {
	return FullProperty(key, "", true)
}

func FullProperty(key string, defaultValue string, isSecret bool)string {
	k8s,k8sFound := os.LookupEnv("IS_KUBERNETES")
	var _val string

	if k8sFound {
		if k8s == "true" || k8s == "True" || k8s == "TRUE"{
			k8sFound = true
		} else {
			k8sFound = false
		}
	}

	if k8sFound {
		config,err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset,err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		if isSecret {
			secret,err := clientset.CoreV1().Secrets("golodex").Get(context.TODO(), "golodex-data", metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			if err != nil {
				_val = ""
			}
			_val = string(secret.Data[key])
		} else {
			config,err := clientset.CoreV1().ConfigMaps("golodex").Get(context.TODO(), "golodex-data", metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			if err != nil {
				_val = ""
			}
			_val = config.Data[key]
		}
	} else {
		_value,found := os.LookupEnv(key)
		log.Println(key + " = " + _value)
		if !found {
			_value = defaultValue
		}
		_val = _value
	}
	return _val
}
