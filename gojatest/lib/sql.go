package lib

import (
	"database/sql"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

type sqlRuntime struct {
	runtime *goja.Runtime
}

type sqlObj struct {
	runtime *goja.Runtime
	db      *sql.DB
}

type rowsObj struct {
	runtime *goja.Runtime
	rows    *sql.Rows
}

func (This *rowsObj) err(call goja.FunctionCall) goja.Value {
	err := This.rows.Err()
	return This.runtime.ToValue(err)
}

func (This *rowsObj) scan(call goja.FunctionCall) goja.Value {
	// 		err = rows.Scan(&id, &username, &realname, &password, &createdAt, &updatedAt)
	// if err != nil {
	// 	panic(This.runtime.NewGoError(err))
	// }
	return nil
}

func (This *rowsObj) next(call goja.FunctionCall) goja.Value {
	r := This.rows.Next()
	return This.runtime.ToValue(r)
}

func (This *rowsObj) close(call goja.FunctionCall) goja.Value {
	err := This.rows.Close()
	if err != nil {
		panic(This.runtime.NewGoError(err))
	}
	return nil
}

func (This *sqlObj) query(call goja.FunctionCall) goja.Value {
	query := call.Argument(0).String()
	rows, err := This.db.Query(query)
	if err != nil {
		panic(This.runtime.NewGoError(err))
	}
	obj := &rowsObj{
		runtime: This.runtime,
		rows:    rows,
	}
	o := This.runtime.NewObject()
	o.Set("close", obj.close)
	o.Set("next", obj.next)
	o.Set("scan", obj.scan)
	o.Set("err", obj.err)
	return o
}

func (This *sqlRuntime) newFunc(call goja.FunctionCall) goja.Value {
	driverName := call.Argument(0).String()
	dataSourceName := call.Argument(1).String()
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(This.runtime.NewGoError(err))
	}
	obj := &sqlObj{
		runtime: This.runtime,
		db:      db,
	}
	o := This.runtime.NewObject()
	o.Set("query", obj.query)
	return o
}

func init() {
	require.RegisterNativeModule("sql", func(runtime *goja.Runtime, module *goja.Object) {
		This := &sqlRuntime{
			runtime: runtime,
		}

		o := module.Get("exports").(*goja.Object)
		o.Set("new", This.newFunc)
	})
}
