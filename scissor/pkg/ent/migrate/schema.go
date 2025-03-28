// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ArticlesColumns holds the columns for the "articles" table.
	ArticlesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint, Increment: true},
		{Name: "title", Type: field.TypeString},
		{Name: "content", Type: field.TypeString, Size: 2147483647},
		{Name: "url", Type: field.TypeString, Unique: true},
		{Name: "author", Type: field.TypeString},
		{Name: "source", Type: field.TypeString},
		{Name: "summary", Type: field.TypeString, Nullable: true},
		{Name: "tags", Type: field.TypeJSON, Nullable: true},
		{Name: "published_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "user_articles", Type: field.TypeInt},
	}
	// ArticlesTable holds the schema information for the "articles" table.
	ArticlesTable = &schema.Table{
		Name:       "articles",
		Columns:    ArticlesColumns,
		PrimaryKey: []*schema.Column{ArticlesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "articles_users_articles",
				Columns:    []*schema.Column{ArticlesColumns[11]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "article_title",
				Unique:  false,
				Columns: []*schema.Column{ArticlesColumns[1]},
			},
			{
				Name:    "article_url",
				Unique:  false,
				Columns: []*schema.Column{ArticlesColumns[3]},
			},
			{
				Name:    "article_author",
				Unique:  false,
				Columns: []*schema.Column{ArticlesColumns[4]},
			},
			{
				Name:    "article_published_at",
				Unique:  false,
				Columns: []*schema.Column{ArticlesColumns[8]},
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString},
		{Name: "email", Type: field.TypeString, Nullable: true},
		{Name: "wx_open_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "nickname", Type: field.TypeString, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_username",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[1]},
			},
			{
				Name:    "user_email",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[3]},
			},
			{
				Name:    "user_wx_open_id",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ArticlesTable,
		UsersTable,
	}
)

func init() {
	ArticlesTable.ForeignKeys[0].RefTable = UsersTable
}
