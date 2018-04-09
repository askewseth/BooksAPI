#!/bin/bash

help() {
	echo "Available commands are:"
	echo " - create: create the pod"
	echo " - expose: expose the ports of the pod after it is created"
	echo " - status: print the status of the pod"
	echo " - delete: delete the pod"
}
	
create() {
	kubectl create -f pod-books-api.yml
	exit $?
}

expose() {
	kubectl port-forward books-api 5555
	exit $?
}

delete() {
	kubectl delete pod books-api
	exit $?
}

status() {
	kubectl get pods books-api | awk '{print $3}' | tail -n1
	exit $?
}

case "$1" in
	expose) expose ;;
	create) create ;;
	delete) delete ;;
	status) status ;;
	*)      help ;;
esac
