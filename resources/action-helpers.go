package resourceprovider

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	echo "github.com/labstack/echo/v4"
	"github.com/ncpa0/hardwire/configuration"
)

func bindFormParams(ctx echo.Context, bodyPtr interface{}) error {
	go (func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("internal error binding form params: ", r)
		}
	})()

	bodyValueElem := reflect.ValueOf(bodyPtr).Elem()
	body := bodyValueElem.Interface()

	if param, err := ctx.FormParams(); err == nil {
		v := reflect.ValueOf(body)
		// iterate over keys of the body struct
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// only bind fields of type string
			if field.Kind() == reflect.String {
				fieldName := reflect.TypeOf(body).Field(i).Name
				paramValue := param.Get(fieldName)
				if paramValue != "" {
					// assign the value to the
					bodyValueElem.Field(i).SetString(paramValue)
				}
			}
		}
		return nil
	} else {
		return err
	}
}

type ActionMetadata struct {
	Resource  string   `json:"resource"`
	Action    string   `json:"action"`
	Method    string   `json:"method"`
	IslandIDs []string `json:"islandIDs"`
}

type ActionsMetadata struct {
	RegisteredActions []ActionMetadata `json:"registeredActions"`
}

func ValidateActionEndpoints() {
	outDir := configuration.Current.HtmlDir
	actionsMetaFilepath := filepath.Join(outDir, "__actions.meta.json")

	fileContent, err := os.ReadFile(actionsMetaFilepath)
	if err != nil {
		panic("Unable to read the actions metadata file")
	}

	// unmarchal the json file
	var actionsMeta ActionsMetadata
	err = json.Unmarshal(fileContent, &actionsMeta)

	if err != nil {
		panic("Actions metadata file is corrupted")
	}

	for _, actionMeta := range actionsMeta.RegisteredActions {
		res, found := ResourceReg.find(actionMeta.Resource)
		if !found {
			panic("Resource referenced by one of the actions doesn't exist: " + actionMeta.Resource)
		}

		found, _ = res.findAction(actionMeta.Method, actionMeta.Action)
		if !found {
			panic("Action used does not exist: " + actionMeta.Method + "/" + actionMeta.Action)
		}
	}
}
