#!/bin/bash
# DO NOT MODIFY THIS SECTION !!!
echo "do migrate"
{{cookiecutter.project_name}} migrate

echo "[Sync] BEGIN ====================="
echo "[Sync] generate definition.yaml"
{{cookiecutter.project_name}} gen_definition_yaml
if [ $? -ne 0 ]
then
	echo "run generate_definition_yaml fail, please run this command on your development env to find out the reason"
	exit 1
fi


if [ -f /app/definition.yaml ]
then
    echo "[Sync] the /app/definition.yaml content:"
	cat /app/definition.yaml
	echo "===================="
fi


echo "[Sync] generate resources.yaml"
{{cookiecutter.project_name}} gen_resources_yaml
if [ $? -ne 0 ]
then
	echo "run generate_resources_yaml fail, please run this command on your development env to find out the reason"
	exit 1
fi

if [ -f /app/resources.yaml ]
then
	echo "[Sync] the /app/resources.yaml content:"
	cat /app/resources.yaml
	echo "===================="
fi

echo "[Sync] sync to apigateway"
{{cookiecutter.project_name}} sync_apigateway
if [ $? -ne 0 ]
then
	echo "run sync_apigateway fail"
	exit 1
fi
echo "[Sync] DONE ====================="
# DO NOT MODIFY THIS SECTION !!! SECTION !!!