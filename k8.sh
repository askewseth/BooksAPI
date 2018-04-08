#!/bin/bash


help() {
	echo "Available commands are:"
	echo " - expose: expose the ports of the pod after it is created"
	echo " - create: create the pod"
	echo " - delete: delete the pod"
}
	
create() {
	kubectl create -f pod-books-api.yml
	exit 0
}

expose() {
	kubectl port-forward books-api 5555
	exit 0
}

delete() {
	kubectl delete pod books-api
	exit 0
}

case "$1" in
	expose) expose ;;
	create) create ;;
	delete) delete ;;
esac

help
