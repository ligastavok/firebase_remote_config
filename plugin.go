package firebase_remote_config

import (
	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"golang.org/x/oauth2/google"

	fb "google.golang.org/api/firebaseremoteconfig/v1"
)

const (
	CHANNEL_NAME_METHOD        = "plugins.flutter.io/firebase_remote_config"
	INSTANCE_METHOD            = "RemoteConfig#instance"
	SET_CONFIG_SETTINGS_METHOD = "RemoteConfig#setConfigSettings"
	FETCH_METHOD               = "RemoteConfig#fetch"
	ACTIVATE_METHOD            = "RemoteConfig#activate"
	SET_DEFAULTS_METHOD        = "RemoteConfig#setDefaults"
	FB_SERVICE_ACCOUNT_FILE    = "fb_service_account.json"
)

type FirebaseRemoteConfigGoFlutterPlugin struct {
	channel *plugin.MethodChannel
}

type FbMessageContent struct {
	Parameters map[string]FbMessageDefaultValue `json:"parameters"`
}

type FbMessageDefaultValue struct {
	DefaultValue FbMessageValue `json:"defaultValue"`
}

type FbMessageValue struct {
	Value interface{} `json:"value"`
}

var _ flutter.Plugin = &FirebaseRemoteConfigGoFlutterPlugin{}

func (p *FirebaseRemoteConfigGoFlutterPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	p.channel = plugin.NewMethodChannel(messenger, CHANNEL_NAME_METHOD, plugin.StandardMethodCodec{})
	p.channel.HandleFunc(INSTANCE_METHOD, handleInstance)
	p.channel.HandleFunc(SET_CONFIG_SETTINGS_METHOD, handleSetConfigSettings)
	p.channel.HandleFunc(FETCH_METHOD, handleFetch)
	p.channel.HandleFunc(ACTIVATE_METHOD, handleActivate)
	p.channel.HandleFunc(SET_DEFAULTS_METHOD, handleSetDefaults)

	return nil
}

// firebase service
func createFirebaseService(ctx context.Context) (*fb.Service, string, error) {
	data, err := ioutil.ReadFile(FB_SERVICE_ACCOUNT_FILE)

	if err != nil {
		fmt.Println("Failed to create firebase service: ", err)
		return nil, "unknown", err
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/firebase.remoteconfig")
	if err != nil {
		fmt.Println("Failed to get config json file: ", err)
		return nil, "unknown", err
	}

	// get project name
	var dataJson map[string]interface{}
	err = json.Unmarshal(data, &dataJson)
	if err != nil {
		fmt.Println("Failed to get config json file: ", err)
		return nil, "unknown", err
	}

	projectName := dataJson["project_id"].(string)
	service, err := fb.New(conf.Client(ctx))

	return service, projectName, err
}

// firebase fetch
func fetchConfig(ctx context.Context) (*fb.RemoteConfig, error) {
	service, projectName, err := createFirebaseService(ctx)

	if err != nil {
		fmt.Println("Failed to create firebase service: ", err)
		return nil, fmt.Errorf("Failed to initialize Firebase service")
	}

	projectID := "projects/" + projectName

	cfg, err := service.Projects.GetRemoteConfig(projectID).Do()
	if err != nil {
		fmt.Println("Failed to call Firebase remote config API: ", err)
		return nil, err
	}

	return cfg, nil
}

func handleSetConfigSettings(arguments interface{}) (reply interface{}, err error) {
	return nil, errors.New("Unimplemented")
}

func handleFetch(arguments interface{}) (reply interface{}, err error) {
	return map[interface{}]interface{}{
		"lastFetchTime":   int64(1596960491970), // stub
		"lastFetchStatus": "success",            // stub
	}, nil
}

// Example of returning result:
// return map[interface{}]interface{} {
// 	"parameters": map[interface{}]interface{} {
// 		"key1": map[interface{}]interface{} {
// 			"source": "remote", // stub
// 			"value": []byte(`{"test_data":[{"activePage":0,"id":1,"label":"some label"}]}`),
// 		},
// 	},
// 	"newConfig": false, // stub
// }, nil
func handleActivate(arguments interface{}) (reply interface{}, err error) {
	// fetch firebase config
	ctx := context.Background()

	cfg, err := fetchConfig(ctx)
	if err != nil {
		// fmt.Println("Failed to call Firebase remote config API: ", err)
		return nil, errors.New("Failed to call Firebase remote config API")
	}

	raw, err := cfg.MarshalJSON()
	if err != nil {
		// fmt.Println("Failed to parse remote config json: ", err)
		return nil, errors.New("Failed to parse remote config json")
	}

	var remoteConfigMap FbMessageContent
	err = json.Unmarshal(raw, &remoteConfigMap)
	if err != nil {
		// fmt.Println("Failed to unmarshal config json: ", err)
		return nil, errors.New("Failed to unmarshal config json")
	}

	var result = make(map[interface{}]interface{})
	var parametersMap = make(map[interface{}]interface{})

	for k := range remoteConfigMap.Parameters {
		parametersMap[k] = map[interface{}]interface{}{
			"source": "remote", // stub
			"value":  []byte(remoteConfigMap.Parameters[k].DefaultValue.Value.(string)),
		}
	}
	result["parameters"] = parametersMap

	return result, nil
}

func handleSetDefaults(arguments interface{}) (reply interface{}, err error) {
	return nil, errors.New("Unimplemented")
}

func handleInstance(arguments interface{}) (reply interface{}, err error) {
	return map[interface{}]interface{}{
		"lastFetchTime":   int64(1596960491970), // stub
		"lastFetchStatus": "success",            // stub
		"inDebugMode":     true,                 // stub
		"parameters": map[interface{}]interface{}{
			"remote": map[interface{}]interface{}{
				"source": "remote",
			},
		},
	}, nil
}
