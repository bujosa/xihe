package utils

import (
	"context"

	"github.com/bujosa/xihe/env"
)

type ContextKey string

const (
	ProjectIdKey     ContextKey = "projectId"
	BucketNameKey    ContextKey = "bucketName"
	SubFolderPathKey ContextKey = "subFolderPath"
	CountryVersionId ContextKey = "countryVersionId"
	CategoryId       ContextKey = "categoryId"
	GeoCode          ContextKey = "geoCode"
	DbUri            ContextKey = "dbUri"
	ProductionApiUrl ContextKey = "productionApiUrl"
	SessionSecret    ContextKey = "sessionSecret"
	CityId           ContextKey = "cityId"
	SpotId           ContextKey = "spotId"
	ErrorCount       ContextKey = "errorCount"
)

func LoadEnvs(ctx *context.Context) {
	projectId, err := env.GetString("GOOGLE_PROJECT_ID")
	if err != nil {
		panic(err)
	}
	bucketName, err := env.GetString("GOOGLE_BUCKET_NAME")
	if err != nil {
		panic(err)
	}

	subFolderPath, err := env.GetString("GOOGLE_SUBFOLDER_PATH")
	if err != nil {
		subFolderPath = ""
	}

	countryVersion, err := env.GetString("COUNTRY_VERSION_ID")
	if err != nil {
		panic(err)
	}

	category, err := env.GetString("CATEGORY_ID")
	if err != nil {
		panic(err)
	}

	value, err := env.GetString("GEOCODE")
	if err != nil {
		panic(err)
	}

	databaseURL, err := env.GetString(DB_URL_ENV_KEY)
	if err != nil {
		panic(err)
	}

	productionURL, err := env.GetString("PRODUCTION_API_URL")
	if err != nil {
		panic(err)
	}
	token, err := env.GetString("SESSION_SECRET")
	if err != nil {
		panic(err)
	}

	city, err := env.GetString("CITY_ID")
	if err != nil {
		panic(err)
	}

	spot, err := env.GetString("SPOT_ID")
	if err != nil {
		panic(err)
	}

	*ctx = context.WithValue(*ctx, ProjectIdKey, projectId)
	*ctx = context.WithValue(*ctx, BucketNameKey, bucketName)
	*ctx = context.WithValue(*ctx, SubFolderPathKey, subFolderPath)
	*ctx = context.WithValue(*ctx, CountryVersionId, countryVersion)
	*ctx = context.WithValue(*ctx, CategoryId, category)
	*ctx = context.WithValue(*ctx, GeoCode, value)
	*ctx = context.WithValue(*ctx, DbUri, databaseURL)
	*ctx = context.WithValue(*ctx, ProductionApiUrl, productionURL)
	*ctx = context.WithValue(*ctx, SessionSecret, token)
	*ctx = context.WithValue(*ctx, CityId, city)
	*ctx = context.WithValue(*ctx, SpotId, spot)
	*ctx = context.WithValue(*ctx, ErrorCount, 0)
}
