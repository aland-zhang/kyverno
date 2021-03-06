<small>*[documentation](/README.md#documentation) / [Writing Policies](/documentation/writing-policies.md) / Configmap Lookup*</small>

# Using ConfigMaps for Variables

There are many cases where the values that are passed into Kyverno policies are dynamic or need to be vary based on the execution environment.

Kyverno supports using Kubernetes [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) to manage variable values outside of a policy definition. 

# Defining ConfigMaps in a Rule Context

To refer to values from a ConfigMap inside any Rule, define a context inside the rule with one or more ConfigMap declarations.

````yaml
  rules:
    - name: example-configmap-lookup
      # added context to define the configmap information which will be referred 
      context:
      # unique name to identify configmap
      - name: dictionary
        configMap: 
          # configmap name - name of the configmap which will be referred
          name: mycmap
          # configmap namepsace - namespace of the configmap which will be referred
          namespace: test
````

Sample ConfigMap Definition

````yaml
apiVersion: v1
data:
  env: production
kind: ConfigMap
metadata:
  name: mycmap
````

# Looking up values

A ConfigMap that is defined in a rule context can be referred to using its unique name within the context. ConfigMap values can be referenced using a JMESPATH style expression `{{<name>.<data>.<key>}}`.

For the example above, we can refer to a ConfigMap value using `{{dictionary.data.env}}`. The variable will be substituted with the value `production` during policy execution.

# Handling Array Values

The ConfigMap value can be an array of string values in JSON format. Kyverno will parse the JSON string to a list of strings, so set operations like In and NotIn can then be applied.

For example, a list of allowed roles can be stored in a ConfigMap, and the Kyverno policy can refer to this list to deny the requests where the role does not match one of the values in the list.

Here are the allowed roles in the ConfigMap:

````yaml
apiVersion: v1
data:
  allowed-roles: "[\"cluster-admin\", \"cluster-operator\", \"tenant-admin\"]"
kind: ConfigMap
metadata:
  name: roles-dictionary
  namespace: test
````

Here is a rule to block a Deployment if the value of annotation `role` is not in the allowed list:

````yaml
spec:
  validationFailureAction: enforce
  rules:
  - name: validate-role-annotation
    context:
      - name: roles-dictionary
        configMap: 
          name: roles-dictionary
          namespace: test
    match:
      resources:
        kinds:
        - Deployment
    preconditions:
    - key: "{{ request.object.metadata.annotations.role }}"
      operator: NotEquals
      value: ""
    validate:
      message: "role {{ request.object.metadata.annotations.role }} is not in the allowed list {{ \"roles-dictionary\".data.\"allowed-roles\" }}"
      deny:
        conditions: 
        - key: "{{ request.object.metadata.annotations.role }}"
          operator: NotIn
          value:  "{{ \"roles-dictionary\".data.\"allowed-roles\" }}"
````



<small>*Read Next >> [Testing Policies](/documentation/testing-policies.md)*</small>
