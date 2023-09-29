# Kanister Visualiser 

The Kanister Visualiser is a user-friendly web application designed to provide a clear and intuitive overview of Kanister's critical components within a Kubernetes cluster. With its sleek and stylish interface, it offers instant insights into Kanister Profiles, Blueprints, and ActionSets. Users can easily identify and explore various Kanister Profiles, including their types and associated cloud storage buckets, thanks to the app's organised display. Furthermore, the visualiser showcases Kanister Blueprints, allowing users to view and comprehend each blueprint's purpose and parameters. 

Perhaps most importantly, the ActionSets section presents a real-time view of ongoing operations, including names, namespaces, and status updates. The app's dark mode feature enhances user experience, and it also provides quick links to Kanister's GitHub repository, website, documentation, and Slack community. 

Overall, the Kanister Visualiser empowers users to monitor and manage their Kanister-based data protection workflows with elegance and ease.

## Getting Started 

On windows you need to run the following to set the KUBECONFIG 

`$env:KUBECONFIG = "C:\Users\micha\.kube\config"`

## Currently working on: 

messing around with creating a dockerfile now 

```
docker run -e KUBECONFIG=/path/to/your/kubeconfig.yaml -p 8080:8080 golang-kanister-app:local

docker run -e KUBECONFIG="C:\Users\micha\.kube\config" -p 8080:8080 golang-kanister-app:local

docker run -v c:\users\micha\.kube\ -v c:\users\micha\.kube\config -e KUBECONFIG=c:\users\micha\.kube\config -p 8080:8080 golang-kanister-app:local
```

To run you should download the source files and use go on your local system to run/build 

```
go run main.go
```

the above command will start the web server on http://localhost:8080/ navigate here and providing you have a KUBECONFIG and valid Kubernetes cluster you will be able to at least see the building blocks of the application, to see the value of the visualiser you will need to continue with the deployment of Kanister and custom resources. 

## Deploy Kanister to your Kubernetes cluster 

`helm repo add kanister https://charts.kanister.io/`

`helm repo update` 

```
helm install kanister --namespace kanister --create-namespace \
  kanister/kanister-operator
```

Check out the custom resources we now have within our Kubernetes cluster 

`kubectl get crds | grep "kanister"`

## Installing tools 

There are two tools including in this script one is `kando` which we won't be using natively and `kanctl`  

`kanctl` simplifies this process by allowing the user to create custom Kanister resources - ActionSets and Profiles, override existing ActionSets and validate profiles.

`curl https://raw.githubusercontent.com/kanisterio/kanister/master/scripts/get.sh | bash`

We will be using `kanctl` to create our example profiles and actionsets. 

## Adding Profiles 

AWS S3 
```
kanctl create profile \
   --namespace kanister \
   --bucket mc-kanister-aws \
   s3compliant \
   --access-key <AWS_ACCESS_KEY> \
   --secret-key <AWS_SECRET_KEY> \
   --region us-east-2
```
Azure 
```
kanctl create profile azure --namespace kanister --storage-account <STORAGE-ACCOUNT> --storage-key <STORAGE-KEY> --bucket <BUCKET-NAME> --region eastus2
```

Google 
```
kanctl create profile gcp --namespace kanister --project-id <PROJECTID>--service-key <SERVICE-KEY> --bucket <BUCKET-NAME> --region us-east1
```

## Deploy your application 

## Adding Blueprints to your cluster 
```
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/csi-snapshot/csi-snapshot-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/postgresql/postgres-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mysql/mysql-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/elasticsearch/elasticsearch-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/aurora-mysql/rds-aurora-snap-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/postgresql/rds-postgres-dump-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/postgresql/rds-postgres-snap-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/cassandra/cassandra-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/cockroachdb/cockroachdb-blueprint.yaml -n kanister 
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/couchbase/couchbase-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/etcd/etcd-in-cluster/k8s/etcd-incluster-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/foundationdb/foundationdb-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/k8ssandra/k8ssandra-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/kafka/adobe-s3-connector/kafka-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/maria/maria-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mongodb/mongo-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mssql/mssql-blueprint.yaml -n kanister
kubectl create -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/redis/redis-blueprint.yaml -n kanister
```
## Create an ActionSet to Protect your application to your profile 
```
kanctl create actionset --action backup --namespace kanister --blueprint mysql-blueprint --statefulset mysql-test/mysql-release --profile kanister/s3-profile-72pvq --secrets mysql=mysql-test/mysql-release
```

```
## Deleting Blueprints 

kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/csi-snapshot/csi-snapshot-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/postgresql/postgres-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mysql/mysql-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/elasticsearch/elasticsearch-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/aurora-mysql/rds-aurora-snap-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/postgresql/rds-postgres-dump-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/aws-rds/postgresql/rds-postgres-snap-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/cassandra/cassandra-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/cockroachdb/cockroachdb-blueprint.yaml -n kanister 
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/couchbase/couchbase-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/etcd/etcd-in-cluster/k8s/etcd-incluster-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/foundationdb/foundationdb-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/k8ssandra/k8ssandra-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/kafka/adobe-s3-connector/kafka-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/maria/maria-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mongodb/mongo-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/mssql/mssql-blueprint.yaml -n kanister
kubectl delete -f https://raw.githubusercontent.com/kanisterio/kanister/master/examples/redis/redis-blueprint.yaml -n kanister

```
