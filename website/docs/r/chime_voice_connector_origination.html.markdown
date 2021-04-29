---
subcategory: "Chime"
layout: "aws"
page_title: "AWS: aws_chime_voice_connector_origination"
description: |-
Enable origination settings to control inbound calling to your SIP infrastructure.
---

# Resource: aws_chime_voice_connector_origination

Enable origination settings to control inbound calling to your SIP infrastructure.

```terraform
resource "aws_chime_voice_connector" "default" {
  name               = "test"
  require_encryption = true
}

resource "aws_chime_voice_connector_origination" "default" {
  disabled            = false
  voice_connector_id  = aws_chime_voice_connector.default.id

  route {
    host = "127.0.0.1"
    port = 8081
    protocol = "TCP"
    priority = 1
    weight = 1
  }

  route {
    host = "127.0.0.2"
    port = 8082
    protocol = "TCP"
    priority = 2
    weight = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `voice_connector_id` - (Required) The Amazon Chime Voice Connector ID.
* `route` - (Required) The call distribution properties defined for your SIP hosts. Valid range: Minimum value of 1. Maximum value of 20.
* `disabled` - (Optional) When origination settings are disabled, inbound calls are not enabled for your Amazon Chime Voice Connector.

### `route`

Origination routes define call distribution properties for your SIP hosts to receive inbound calls using your Amazon Chime Voice Connector. Limit: Ten origination routes for each Amazon Chime Voice Connector.

* `host` - The FQDN or IP address to contact for origination traffic.
* `port` - The designated origination route port. Defaults to `5060`.
* `priority` - The priority associated with the host, with 1 being the highest priority. Higher priority hosts are attempted first.
* `protocol` - The protocol to use for the origination route. Encryption-enabled Amazon Chime Voice Connectors use TCP protocol by default.
* `weight` - The weight associated with the host. If hosts are equal in priority, calls are redistributed among them based on their relative weight.

## Import

Configuration Recorder can be imported using the name, e.g.

```
$ terraform import aws_chime_voice_connector_origination.default example
```