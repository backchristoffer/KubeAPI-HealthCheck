### OCP kube-apiserver healthcheck

#### Checks if kube-apiserver is UP
Use .env with the following variable
- BEARER_TOKEN  , using prometheus sa token
- PROM_URL  . Example: https://prometheus-k8s-openshift-monitoring.apps.<...>/api/v1/query

#### Create prometheus-sa if you do not have one
~~~
$ oc create serviceaccount prometheus-sa -n openshift-monitoring
$ oc create token prometheus-sa -n openshift-monitoring
$ oc adm policy add-cluster-role-to-user view -z prometheus-sa -n openshift-monitoring
~~~