# firebase_remote_config

This Go package implements the host-side of the Flutter [firebase_remote_config](https://pub.dev/packages/firebase_remote_config) plugin.

## Usage

1) Import as:

```go
import "github.com/jWinterDay/firebase_remote_config"
```

Then add the following option to your go-flutter [application options](https://github.com/go-flutter-desktop/go-flutter/wiki/Plugin-info):

```go
flutter.AddPlugin(&firebase_remote_config.FirebaseRemoteConfigGoFlutterPlugin{}),
```

2) You must download firebase service account [service-account-file.json](https://firebase.google.com/docs/admin/setup)

3) Create empty file in the root of your Flutter application with name fb_service_account.json

4) Copy-paste content of your service-account-file.json to nearly created file fb_service_account.json

# Working

The plugin is still under development and I need your help as pull-requests

# TODO

- save data local (local storage for example)

# Issues

https://github.com/jWinterDay/firebase_remote_config/issues