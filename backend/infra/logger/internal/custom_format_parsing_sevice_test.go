package internal

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParsingCustomFormat(t *testing.T) {
	Convey("Given a custom format as plog", t, func() {
		format := "${T}\t${L}\t${H}\t${N}\t" +
			"TransactionID=${TransactionID:null}\tInstanceID=${InstanceID:null}\t" +
			"[ObjectID=${ObjectID:null},ObjectType=${ObjectType:null}]\t" +
			"${M} ${E}\t[${C_F:}]:[${C_L}]\n"

		Convey("When parse custom format", func() {
			parsedFormat, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then parsed format should be expect", func() {
				expect := CustomFormat([]CustomFormatItem{
					{
						Field:        "T",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "\t",
					},
					{
						Field:        "L",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "\t",
					},
					{
						Field:        "H",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "\t",
					},
					{
						Field:        "N",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "\tTransactionID=",
					},
					{
						Field:        "TransactionID",
						DefaultValue: "null",
					},
					{
						Field:        "",
						DefaultValue: "\tInstanceID=",
					},
					{
						Field:        "InstanceID",
						DefaultValue: "null",
					},
					{
						Field:        "",
						DefaultValue: "\t[ObjectID=",
					},
					{
						Field:        "ObjectID",
						DefaultValue: "null",
					},
					{
						Field:        "",
						DefaultValue: ",ObjectType=",
					},
					{
						Field:        "ObjectType",
						DefaultValue: "null",
					},
					{
						Field:        "",
						DefaultValue: "]\t",
					},
					{
						Field:        "M",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: " ",
					},
					{
						Field:        "E",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "\t[",
					},
					{
						Field:        "C_F",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "]:[",
					},
					{
						Field:        "C_L",
						DefaultValue: "",
					},
					{
						Field:        "",
						DefaultValue: "]\n",
					},
				})
				So(err, ShouldBeNil)
				So(parsedFormat, ShouldResemble, expect)
			})
		})
	})

	Convey("Given a empty custom format", t, func() {
		format := ""

		Convey("When parse custom format", func() {
			parsedFormat, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then parsed format should be expect", func() {
				expect := CustomFormat(nil)
				So(err, ShouldBeNil)
				So(parsedFormat, ShouldResemble, expect)
			})
		})
	})

	Convey("Given a custom format only has one field", t, func() {
		format := "${XX:xx}"

		Convey("When parse custom format", func() {
			parsedFormat, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then parsed format should be expect", func() {
				expect := CustomFormat([]CustomFormatItem{
					{Field: "XX", DefaultValue: "xx"},
				})
				So(err, ShouldBeNil)
				So(parsedFormat, ShouldResemble, expect)
			})
		})
	})

	Convey("Given a custom format without fields", t, func() {
		format := "xxxxx:xx$xx"

		Convey("When parse custom format", func() {
			parsedFormat, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then parsed format should be expect", func() {
				expect := CustomFormat([]CustomFormatItem{
					{Field: "", DefaultValue: "xxxxx:xx$xx"},
				})
				So(err, ShouldBeNil)
				So(parsedFormat, ShouldResemble, expect)
			})
		})
	})

	Convey("Given a abnormal custom format without the last field ending token", t, func() {
		format := "${T}\t${L}\t${H}\t${N}\t" +
			"TransactionID=${TransactionID:null}\tInstanceID=${InstanceID:null}\t" +
			"[ObjectID=${ObjectID:null},ObjectType=${ObjectType:null}]\t" +
			"${M} ${E}\t[${C_F:}]:[${C_L]\n"

		Convey("When parse custom format", func() {
			_, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then is abnormal and tips: without field ending token", func() {
				So(err.Error(), ShouldContainSubstring, "without field ending token")
			})
		})
	})

	Convey("Given a abnormal custom format without the middle field ending token", t, func() {
		format := "${T}\t${L}\t${H}\t${N}\t" +
			"TransactionID=${TransactionID:null}\tInstanceID=${InstanceID:null}\t" +
			"[ObjectID=${ObjectID:null,ObjectType=${ObjectType:null}]\t" +
			"${M} ${E}\t[${C_F:}]:[${C_L}]\n"

		Convey("When parse custom format", func() {
			_, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then is abnormal and tips: without field ending token", func() {
				So(err.Error(), ShouldContainSubstring, "without field ending token")
			})
		})
	})

	Convey("Given a abnormal custom format without field name", t, func() {
		format := "${:xx}"

		Convey("When parse custom format", func() {
			_, err := GetCustomFormatParsingService().Parse(format)

			Convey("Then is abnormal and tips: without field name", func() {
				So(err.Error(), ShouldContainSubstring, "without field name")
			})
		})
	})
}
