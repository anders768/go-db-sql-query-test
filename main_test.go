package main

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // откладываем закрытие БД

	clientID := 1 // ожидаемое значение, т.е. эталонное; по нему в качестве идентификатора мы делаем запрос

	// напиши тест здесь
	cl, err := selectClient(db, clientID) // получаем объект по идентификатору из clientID
	// проверяем, что функция selectClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции selectClient(): %v", err)

	// проверяем, что ID полученного объекта совпадает со значением в clientID
	require.Equal(t, clientID, cl.ID, "Полученный идентификатор %d не совпадает с ожидаемым %d", cl.ID, clientID)
	// проверяем, что поля полученного объекта - не пустые
	require.NotEmpty(t, cl.FIO, "Поля полученного объекта не должны быть пустыми")
	require.NotEmpty(t, cl.Login, "Поля полученного объекта не должны быть пустыми")
	require.NotEmpty(t, cl.Birthday, "Поля полученного объекта не должны быть пустыми")
	require.NotEmpty(t, cl.Email, "Поля полученного объекта не должны быть пустыми")
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // откладываем закрытие БД

	clientID := -1 // ожидаемое значение, т.е. эталонное; но такого ID не может быть в БД (т.к. < 0)

	// напиши тест здесь
	cl, err := selectClient(db, clientID) // получаем объект по идентификатору из clientID
	// проверяем, что функция selectClient вернула ошибку
	require.Error(t, err, "Функция selectClient() должна была вернуть ошибку")
	// проверяем, то функция selectClient вернула ошибку sql.ErrNoRows
	require.ErrorIs(t, err, sql.ErrNoRows, "Функция selectClient должна была вернуть ошибку sql.ErrNoRows, а не %v", err)
	// проверяем, что поля полученного объекта - пустые
	require.Empty(t, cl.FIO, "Поля полученного объекта должны быть пустыми")
	require.Empty(t, cl.Login, "Поля полученного объекта должны быть пустыми")
	require.Empty(t, cl.Birthday, "Поля полученного объекта должны быть пустыми")
	require.Empty(t, cl.Email, "Поля полученного объекта должны быть пустыми")
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // откладываем закрытие БД

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
	cl.ID, err = insertClient(db, cl) // вставляем строку, получаем её идентификатор
	if err != nil {
		log.Fatal(err)
	}

	// проверяем, что полученный идентификатор - не пустой
	require.NotEmpty(t, cl.ID, "Полученный идентификатор не должен быть пустым")
	// проверяем, что функция insertClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции insertClient(): %v", err)

	cl2, err := selectClient(db, cl.ID) // получаем только что вставленный объект по его полученному идентификатору
	// проверяем, что функция selectClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции selectClient(): %v", err)
	// проверяем равенство полей вставленного и полученного объектов
	require.Equal(t, cl, cl2, "Поля вставленного и полученного объектов должны быть равны") // через сравнение структур
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // откладываем закрытие БД

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
	cl.ID, err = insertClient(db, cl) // вставляем строку, получаем её идентификатор
	if err != nil {
		log.Fatal(err)
	}
	// проверяем, что полученный идентификатор - не пустой
	require.NotEmpty(t, cl.ID, "Полученный идентификатор не должен быть пустым")
	// проверяем, что функция insertClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции insertClient(): %v", err)

	cl, err = selectClient(db, cl.ID) // Получаем только что вставленный объект по его полученному идентификатору;
	// забиваем старые значения полей переменной cl новыми.
	// проверяем, что функция selectClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции selectClient(): %v", err)

	err = deleteClient(db, cl.ID)
	// проверяем, что функция deleteClient не вернула ошибку
	require.NoError(t, err, "Ошибка функции deleteClient(): %v", err)

	cl, err = selectClient(db, cl.ID) // Пытаемся получить только что удалённый объект по его сохранённому идентификатору;
	// видимо, забиваем старые значения полей переменной cl новыми.
	// проверяем, что функция selectClient вернула ошибку
	require.Error(t, err, "Функция selectClient() должна была вернуть ошибку")
	// проверяем, то функция selectClient вернула ошибку sql.ErrNoRows
	require.ErrorIs(t, err, sql.ErrNoRows, "Функция selectClient должна была вернуть ошибку sql.ErrNoRows, а не %v", err)
}
