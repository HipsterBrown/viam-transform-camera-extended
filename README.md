# `transform-camera-extended` module

This [module](https://docs.viam.com/registry/modular-resources/) implements the [`rdk:component:camera` API](https://docs.viam.com/appendix/apis/components/camera/) applying filters and transformations on images from a source camera component.

## Requirements

This module assumes you have an existing camera component, like a [webcam](https://docs.viam.com/operate/reference/components/camera/webcam/) or [CSI camera](https://docs.viam.comhttps//github.com/viam-modules/csi-camera/).

## Configure your pi-ai-camera vision service

Navigate to the [**CONFIGURE** tab](https://docs.viam.com/configure/) of your [machine](https://docs.viam.com/fleet/machines/) in the [Viam app](https://app.viam.com/).
[Add `hipsterbrown:transform-camera-extended` to your machine](https://docs.viam.com/configure/#components).

### Attributes

The following attributes are available for `hipsterbrown:camera:transform-extended` camera component:

| Name    | Type   | Required?    | Default | Description |
| ------- | ------ | ------------ | ------- | ----------- |
| `source` | string | Required | N/A  | Name of the camera component to use when applying filters & transformations  |
| `pipeline` | array of objects * | Optional     | [] | A list of filters or transformations to apply to source camera image |

* See below for description of each pipeline transform

**Filters & Transformations**

Each object in the list of `pipeline` filters & transformations has a required `"type"` field and optional `"params"` field to pass arguments to the underlying image manipulation function.
See the [`imaging` documentation](https://github.com/disintegration/imaging?tab=readme-ov-file#usage-examples) for how each filter & transformation affects the image.

**`resize`**

Params:
- width: number
- height: number

**`brightness`**

Params:
- amount: float (between -100 and 100)

**`gamma`**

Params:
- amount: float (between 0.0 and 2.0)

**`contrast`**

Params:
- amount: float (between -100 and 100)

**`blur`**

Params:
- sigma: float (between 0.0 and 2.0)

**`sharpen`**

Params:
- sigma: float (between 0.0 and 2.0)

**`grayscale`**

**`fliph`**

**`flipv`**

**`rotate`**

Params:
- angle: float (between -360 and 360)

### Example configuration

```json
{
    "source": "my-camera",
    "pipeline": [
        {
            "type": "brightness",
            "params": {
                "amount": -20
            }
        },
        {
            "type": "grayscale"
        }
    ]
}
```
