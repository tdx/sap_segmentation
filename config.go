package sap_segmentation

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	StorerArgs DbConnection
	LoaderArgs ConnArgs
	Log        Logger
}

// DB_HOST - IP адрес DB сервера, значение по-умолч. "127.0.0.1"
// DB_PORT - TCP порт DB сервера, значение по-умолч. "5432"
// DB_NAME - Название DB, значение по-умолч. "mesh_group"
// DB_USER - Имя пользователя DB, , значение по-умолч. "postgres"
// DB_PASSWORD - Пароль пользователя DB, значение по-умолч. "postgres"

type DbConnection struct {
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"5432"`
	DbName   string `default:"mesh_group" envconfig:"DB_NAME"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
}

// CONN_URI - Адрес для подключения к внешнему API, значение по-умолч. "http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation"
// CONN_AUTH_LOGIN_PWD - Аутентификационные логин и пароль, значение по-умолч. "4Dfddf5:jKlljHGH"
// CONN_USER_AGENT - Юзер агент для подключения к САПу, значение по-умолч. "spacecount-test"
// CONN_TIMEOUT - Время таймаута подключения к внешнему API,в секундах, значение по-умолч. "5"
// CONN_INTERVAL - Задержка между получением очередной пачки данных из САП, значение по-умолч. 1500мс
// IMPORT_BATCH_SIZE - Размер пачки данных для получения при каждом запросе к внешнему API, значение по-умолч. "50"

type ConnArgs struct {
	Uri             string        `default:"http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation"`
	AuthLoginPwd    string        `default:"4Dfddf5:jKlljHGH" split_words:"true"`
	UserAgent       string        `default:"spacecount-test" split_words:"true"`
	ConnectTimeout  time.Duration `default:"5s" envconfig:"CONN_TIMEOUT"`
	Interval        time.Duration `default:"1500ms"`
	ImportBatchSize int           `default:"50" envconfig:"IMPORT_BATCH_SIZE"`
}

// LOG_CLEANUP_MAX_AGE - Время, после которого удаляются старые логи, в днях, значение по-умолч. "7"

type Logger struct {
	CleanupMaxAgeDays int `default:"7" envconfig:"LOG_CLEANUP_MAX_AGE"`
}

func LoadConfig() (*Config, error) {
	var db DbConnection
	if err := envconfig.Process("DB", &db); err != nil {
		return nil, err
	}

	var conn ConnArgs
	if err := envconfig.Process("CONN", &conn); err != nil {
		return nil, err
	}

	var lg Logger
	if err := envconfig.Process("LOG", &lg); err != nil {
		return nil, err
	}

	return &Config{
		StorerArgs: db,
		LoaderArgs: conn,
		Log:        lg}, nil
}
