# Azure Virtual Network Flow Logs Receiver

This receiver reads Azure virtual network flow logs from [Azure Blob Storage](https://azure.microsoft.com/en-us/products/storage/blobs/).

## Configuration

The following settings are required:

- `event_hub:`
  `  endpoint:` (no default): Azure Event Hub endpoint triggering on the `Blob Create` event 

The following settings can be optionally configured:

- `auth` (default = connection_string): Specifies the used authentication method. Supported values are `connection_string`, `service_principal`, `default`.
- `cloud` (default = "AzureCloud"): Defines which Azure Cloud to use when using the `service_principal` authentication method. Value is either `AzureCloud` or `AzureUSGovernment`.
- `logs:`
  `  container_name:` (default = "insights-logs-flowlogflowevent"): Name of the blob container with the logs

Authenticating using a connection string requires configuration of the following additional setting:

- `connection_string:` Azure Blob Storage connection key, which can be found in the Azure Blob Storage resource on the Azure Portal.

Authenticating using service principal requires configuration of the following additional settings:

- `service_principal:`
  `  tenant_id`
  `  client_id`
  `  client_secret`
- `storage_account_url:` Azure Storage Account url

The service principal method also requires the [Storage Blob Data Contributor](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles/storage#storage-blob-data-contributor) role on the logs containers.

### Example configurations

Using connection string for authentication:

```yaml
receivers:
  azureblob:
    connection_string: DefaultEndpointsProtocol=https;AccountName=accountName;AccountKey=+idLkHYcL0MUWIKYHm2j4Q==;EndpointSuffix=core.windows.net
    event_hub:
      endpoint: Endpoint=sb://oteldata.servicebus.windows.net/;SharedAccessKeyName=otelhubbpollicy;SharedAccessKey=mPJVubIK5dJ6mLfZo1ucsdkLysLSQ6N7kddvsIcmoEs=;EntityPath=otellhub
```

Using service principal for authentication:

```yaml
receivers:
  azureblob:
    auth: service_principal
    service_principal:
      tenant_id: "${tenant_id}"
      client_id: "${client_id}"
      client_secret: "${env:CLIENT_SECRET}"
    storage_account_url: https://accountName.blob.core.windows.net
    event_hub:
      endpoint: Endpoint=sb://oteldata.servicebus.windows.net/;SharedAccessKeyName=otelhubbpollicy;SharedAccessKey=mPJVubIK5dJ6mLfZo1ucsdkLysLSQ6N7kddvsIcmoEs=;EntityPath=otellhub
```

The receiver subscribes [on the events](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-event-overview) published by Azure Blob Storage and handled by Azure Event Hub. When it receives `Blob Create` event, it reads the logs from a corresponding blob and deletes it after processing.

