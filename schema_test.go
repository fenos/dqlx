package deku_test

import (
	"github.com/fenos/deku"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TypeBuilder(t *testing.T) {
	t.Run("generate a simple type", func(t *testing.T) {
		schema := deku.NewSchema()
		carType := schema.
			Type("Car").
			String("name").
			String("brand")

		writer := deku.NewWriter()
		writer.AddLine("type Car {")
		writer.AddIndentedLine("Car.name")
		writer.AddIndentedLine("Car.brand")
		writer.AddLine("}")
		expected := writer.ToString()

		require.Equal(t, expected, carType.ToString())
	})

	t.Run("generate a simple type without prefix", func(t *testing.T) {
		schema := deku.NewSchema()
		carType := schema.
			Type("Car", deku.WithPrefix(false)).
			String("name").
			String("brand")

		writer := deku.NewWriter()
		writer.AddLine("type Car {")
		writer.AddIndentedLine("name")
		writer.AddIndentedLine("brand")
		writer.AddLine("}")
		expected := writer.ToString()

		require.Equal(t, expected, carType.ToString())
	})

	t.Run("can register the same predicate", func(t *testing.T) {
		schema := deku.NewSchema()

		carType := schema.
			Type("Car", deku.WithPrefix(false)).
			String("name")

		busType := schema.
			Type("Bus", deku.WithPrefix(false)).
			String("name")

		carWriter := deku.NewWriter()
		carWriter.AddLine("type Car {")
		carWriter.AddIndentedLine("name")
		carWriter.AddLine("}")

		busWriter := deku.NewWriter()
		busWriter.AddLine("type Bus {")
		busWriter.AddIndentedLine("name")
		busWriter.AddLine("}")

		require.Equal(t, carWriter.ToString(), carType.ToString())
		require.Equal(t, busWriter.ToString(), busType.ToString())
	})

	t.Run("cannot register the same field twice on the same type", func(t *testing.T) {
		schema := deku.NewSchema()
		require.Panicsf(t, func() {
			schema.
				Type("Car").
				String("name").
				String("name")

		}, "predicate '%s' already registered", "name")
	})

	t.Run("cannot register the same type twice", func(t *testing.T) {
		schema := deku.NewSchema()
		require.Panicsf(t, func() {
			schema.
				Type("Car").
				String("name")

			schema.
				Type("Car").
				String("name")
		}, "type 'Car' already registered")
	})
}

func Test_PredicateBuilder(t *testing.T) {

	t.Run("string indexes", func(t *testing.T) {
		schema := deku.NewSchema()

		schema.
			Type("Car", deku.WithPrefix(false)).
			String("name").IndexExact().
			String("brand").IndexHash().
			String("description").IndexFulltext().IndexFulltext(). // it should not duplicate index
			String("manufacture").IndexTerm().
			String("trigram").IndexTrigram().
			String("multiple").IndexTrigram().IndexExact().IndexHash().
			DateTime("created_at").IndexYear().IndexMonth().IndexDay().IndexHour().IndexHour()

		writer := deku.NewWriter()
		writer.AddLine("name:string @index(exact) .")
		writer.AddLine("brand:string @index(hash) .")
		writer.AddLine("description:string @index(fulltext) .")
		writer.AddLine("manufacture:string @index(term) .")
		writer.AddLine("trigram:string @index(trigram) .")
		writer.AddLine("multiple:string @index(trigram,exact,hash) .")
		writer.AddLine("created_at:dateTime @index(year,month,day,hour) .")

		require.Equal(t, writer.ToString(), schema.PredicatesToString())
	})

	t.Run("normal index", func(t *testing.T) {
		schema := deku.NewSchema()

		schema.
			Type("Car", deku.WithPrefix(false)).
			UID("id").Index().
			Int("km").Index().
			Float("float").Index().Index(). // It should not duplicate the index
			Bool("bool").Index().
			Geo("geo").Index().
			Password("password").Index()

		writer := deku.NewWriter()
		writer.AddLine("id:uid @index() .")
		writer.AddLine("km:int @index() .")
		writer.AddLine("float:float @index() .")
		writer.AddLine("bool:bool @index() .")
		writer.AddLine("geo:geo @index() .")
		writer.AddLine("password:password @index() .")

		require.Equal(t, writer.ToString(), schema.PredicatesToString())
	})

	t.Run("upsert / lang / list index", func(t *testing.T) {
		schema := deku.NewSchema()

		schema.
			Type("Car", deku.WithPrefix(false)).
			Int("km").Upsert().
			String("brand").Lang().
			Int("rates").List()

		writer := deku.NewWriter()
		writer.AddLine("km:int @upsert() .")
		writer.AddLine("brand:string @lang() .")
		writer.AddLine("rates:[int] .")

		require.Equal(t, writer.ToString(), schema.PredicatesToString())
	})

	t.Run("reverse relation", func(t *testing.T) {
		schema := deku.NewSchema()

		carType := schema.
			Type("Car", deku.WithPrefix(false)).
			UID("id").
			Int("shops").Reverse()

		writer := deku.NewWriter()
		writer.AddLine("id:uid .")
		writer.AddLine("shops:int @reverse() .")

		require.Equal(t, writer.ToString(), schema.PredicatesToString())

		writerType := deku.NewWriter()
		writerType.AddLine("type Car {")
		writerType.AddIndentedLine("id")
		writerType.AddIndentedLine("<~shops>")
		writerType.AddLine("}")

		require.Equal(t, writerType.ToString(), carType.ToString())
	})
}

func Test_FullSchema(t *testing.T) {
	schema := deku.NewSchema()

	schema.Type("Car").
		UID("id").
		String("name").IndexTerm().
		String("brand").
		Int("hp")

	schema.Type("Shop").
		UID("id").
		String("name").IndexTerm().
		String("street").IndexFulltext().
		Int("branches").Index()


	writer := deku.NewWriter()
	writer.AddLine("Car.id:uid .")
	writer.AddLine("Car.name:string @index(term) .")
	writer.AddLine("Car.brand:string .")
	writer.AddLine("Car.hp:int .")

	writer.AddLine("Shop.id:uid .")
	writer.AddLine("Shop.name:string @index(term) .")
	writer.AddLine("Shop.street:string @index(fulltext) .")
	writer.AddLine("Shop.branches:int @index() .")

	writer.AddLine("type Car {")
	writer.AddIndentedLine("Car.id")
	writer.AddIndentedLine("Car.name")
	writer.AddIndentedLine("Car.brand")
	writer.AddIndentedLine("Car.hp")
	writer.AddLine("}")

	writer.AddLine("type Shop {")
	writer.AddIndentedLine("Shop.id")
	writer.AddIndentedLine("Shop.name")
	writer.AddIndentedLine("Shop.street")
	writer.AddIndentedLine("Shop.branches")
	writer.AddLine("}")

	require.Equal(t, writer.ToString(), schema.ToString())
}