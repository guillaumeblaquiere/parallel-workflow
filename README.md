# Overview
This repository present the different way to contact an Oracle Database with GCP serverless Product.

AppEngine standard and flex, Cloud Run and function are used. Except for AppEngine8, the 4 products are usable with the same source code.

Think to customize the configuration files with your values.

# Deployment

The application is composed of 3 parts

* The Cloud Run app that we will use to wait awhile. Its name is **sleepy-app**
* The Workflow to run the Cloud Run app. Its name is **run-long-process**
* The Workflow to run in parallel several workflows. Its name is **parallel-executor**

## Cloud Run sleepy app

It's a very simple app that take a query parameter `w`. This parameter, an integer, is the number of second to wait
before answering. It simulates long-running process such as Dataflow execution or BigQuery query.

To deploy it, we will use buildpack.io
```
gcloud beta run deploy --source=. --region=us-central1 --platform=managed --allow-unauthenticated sleepy-app
```

Get the service URL at the end

## Run-long-process workflow

This workflow calls the sleepy app. The first step is to update the URL in the yaml file.

* Replace the `<SLEEPY-APP URL>` placeholder by the URL of your Cloud Run sleepy-app got at the previous step. 

Then deploy it

```
gcloud workflows deploy --source=workflow/run-long-process.yaml --location=us-central1 run-long-process 
```

And test it

```
gcloud workflows execute --data='{"wait":5}' run-long-process 
```

Get the latest line of the output, wait 5 second and paste the command. Check the `state` of the workflow and
the time duration between the `startTime` and the `endTime`

## Parallel-executor workflow

This workflow has only hardcoded calls. Start by deploying it

```shell
gcloud workflows deploy --source=workflow/parallel-executor.yaml --location=us-central1 parallel-executor
```

And test it *(no data in entry this time, only for the test)*

```
gcloud workflows execute parallel-executor
```

Get the latest line of the output, wait 15 second and paste the command. Check the `state` of the workflow and
the time duration between the `startTime` and the `endTime`

*Feb 2021: the performance aren't so good, but there is a beginning of parallelization. It's yet experimental!*

# Advanced workflow

To bypass the current limitation, you can use these advance workflows.

## Custom API call

If you want to call different URL in the execution map, you can create a wrapper workflow for this. 
*It's simply an example, you can set the HTTP verb and others part (body) also dynamic and pass in the body*

**Remarks**

If you want to omit some parameter, you need to test them before and set a value (null or not). The `prepareQuery` step
is a good example on this

### Deploy the callable workflow

```shell
gcloud workflows deploy --source=workflow/advanced/custom-api-call.yaml custom-api-call
```
And test it
```shell
gcloud workflows execute --data='{"url":"https://sleepy-app-fqffbf2xsq-uc.a.run.app","query":{"w":5}}' custom-api-call

#Validate that the query param omission isn't a problem
gcloud workflows execute --data='{"url":"https://sleepy-app-fqffbf2xsq-uc.a.run.app"}' custom-api-call
```

### Deploy the execution map for custom API calls

Example of use, with static call

* Replace the `<SLEEPY-APP URL>` placeholder by the URL of your Cloud Run sleepy-app got at the first step.

```shell
gcloud workflows deploy --source=workflow/advanced/parallel-executor-custom-api-call.yaml parallel-executor-custom-api-call
```

Test it
```shell
gcloud workflows execute parallel-executor-custom-api-call
```

## Custom Workflow invocation

To have the capacity to invoke different workflow ID in the execution map, you need to wrap the call to a workflow
in another workflow. It's not very clean, but it's the current workaround to gain flexibility.

### Deploy the callable workflow

```shell
gcloud workflows deploy --source=workflow/advanced/custom-workflow.yaml custom-workflow
```

And test it

```shell
gcloud workflows execute --data='{"workflow":"run-long-process","argument":{"wait":5}}' custom-workflow
```

### Deploy the execution map for custom workflows

Example of use, with static call.

* Replace the `<SLEEPY-APP URL>` placeholder by the URL of your Cloud Run sleepy-app got at the first step.

```shell
gcloud workflows deploy --source=workflow/advanced/parallel-executor-custom-workflow.yaml parallel-executor-custom-workflow
```

Test it
```shell
gcloud workflows execute parallel-executor-custom-workflow
```

# License

This library is licensed under Apache 2.0. Full license text is available in
[LICENSE](https://github.com/guillaumeblaquiere/parallel-workflow/tree/master/LICENSE).