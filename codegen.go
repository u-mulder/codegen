package codegen

import (
	"errors"
	"fmt"
)

const (
	// DefaultLineBreak is a typical
	// line break in php scripts
	DefaultLineBreak = "\n"
	// DefaultIndent is a typical
	// PSR-indentation in php scripts
	DefaultIndent = "    "
	// DefaultArrayLine is template for a row
	// of code used in array value declaration
	DefaultArrayLine = "%s'%s' => %s,%s"
)

// ErrGetSnippet ...
var ErrGetSnippet = errors.New("Snippet with such key not found")

// ErrGetGenerator ...
var ErrGetGenerator = errors.New("Generator with such name not found")

// Codegen represents bunch of settings required for code generation
type Codegen struct {
	LineBreak   string
	Indent      string
	phpSnippets map[string]string
	codegens    map[string]func(c *Codegen) string
}

// New creates a new instance of Codegen type
func New() (*Codegen, error) {
	c := &Codegen{
		LineBreak:   DefaultLineBreak,
		Indent:      DefaultIndent,
		phpSnippets: make(map[string]string),
		codegens:    make(map[string]func(c *Codegen) string),
	}

	return c, nil
}

func (c *Codegen) setLineBreak(lb string) {
	c.LineBreak = lb
}

func (c *Codegen) setIndent(indent string) {
	c.Indent = indent
}

// AddSnippet adds new line(s) of code which can later be used in generation procedure
func (c *Codegen) AddSnippet(key, value string) {
	c.phpSnippets[key] = value
}

// GetSnippet returns previously added line(s) of code by it's key or error
func (c *Codegen) GetSnippet(key string) (string, error) {
	var r string
	var err error
	var ok bool

	if r, ok = c.phpSnippets[key]; !ok {
		err = ErrGetSnippet
	}

	return r, err
}

// RegisterGenerator saves new generator which later can be called
func (c *Codegen) RegisterGenerator(genName string, f func(c *Codegen) string) {
	c.codegens[genName] = f
}

// Generate method runs some code and returns string with php-codes or error
func (c *Codegen) Generate(genName string) (string, error) {
	var resultStr string
	var err error

	if g, ok := c.codegens[genName]; ok {
		resultStr = g(c)
	} else {
		err = ErrGetGenerator
	}

	return resultStr, err
}

// EncloseInSingleQuotes does what it's name says)
func EncloseInSingleQuotes(str string) string {
	return "'" + str + "'"
}

// AddDefaultSnippets add some default lines of code
func (c *Codegen) AddDefaultSnippets() {
	/* common props */
	c.AddSnippet("php", "<?php")
	c.AddSnippet("path", "$_SERVER['DOCUMENT_ROOT'] = dirname(__FILE__);")
	c.AddSnippet("bxheader", "require_once($_SERVER['DOCUMENT_ROOT'] . '/bitrix/modules/main/include/prolog_before.php');")

	/* specially for iblock properties */
	c.AddSnippet("iblock", "\\Bitrix\\Main\\Loader::includeModule('iblock');")
	c.AddSnippet("ibp_obj", "$ibp = new CIBlockProperty();")
	c.AddSnippet("iblock_prop_st", "$prop = array(")
	c.AddSnippet("iblock_prop_en", ");")
	c.AddSnippet("iblock_prop_run", "$r = $ibp->add($prop);")
	c.AddSnippet("iblock_prop_check", "if ($r) { "+c.LineBreak+c.Indent+"echo 'Added prop with ID: ' . $r . PHP_EOL; "+c.LineBreak+"} else { "+c.LineBreak+c.Indent+"echo 'Error adding prop: ' . $ibp->LAST_ERROR  . PHP_EOL; "+c.LineBreak+"}")

	/* specially for user fields */
	c.AddSnippet("uf_obj", "$ufo = new CUserTypeEntity;")
	c.AddSnippet("uf_data_st", "$uf = array(")
	c.AddSnippet("uf_data_en", ");")
	c.AddSnippet("uf_data_run", "$r = $ufo->add($uf);")
	c.AddSnippet("uf_data_check", "if ($r) { "+c.LineBreak+c.Indent+"echo 'Added UserField with ID: ' . $r . PHP_EOL; "+c.LineBreak+"} else { "+c.LineBreak+c.Indent+"echo 'Error adding UserField: ' . $ufo->LAST_ERROR  . PHP_EOL; "+c.LineBreak+"}")
	// TODO - no LAST_ERROR here, catch exception?

	/* specially for mail event/template */
	c.AddSnippet("mevent_obj", "$meo = new CEventType;")
	c.AddSnippet("mevent_data_st", "$me = array(")
	c.AddSnippet("mevent_data_en", ");")
	c.AddSnippet("mevent_run", "$r = $meo->add($me);")
	c.AddSnippet("mevent_run_succ_wo_mm", "if ($r) { "+c.LineBreak+c.Indent+"echo 'Added MailEvent with ID: ' . $r . PHP_EOL; "+c.LineBreak+"} else { "+c.LineBreak+c.Indent+"echo 'Error adding MailEvent: ' . $meo->LAST_ERROR  . PHP_EOL; "+c.LineBreak+"}")
	c.AddSnippet("mevent_run_check", "if ($r) {")
	c.AddSnippet("mevent_run_check_else", "} else {"+c.LineBreak+c.Indent+"echo 'Error adding MailEvent: ' . $meo->LAST_ERROR  . PHP_EOL; "+c.LineBreak+"}")
	c.AddSnippet("mevent_succ", c.Indent+"echo 'Added MailEvent with ID: ' . $r . PHP_EOL;"+c.LineBreak)

	c.AddSnippet("mtpl_warn", c.Indent+"// TODO - set proper LID for template!")
	c.AddSnippet("mtpl_obj", "$mmo = new CEventMessage;")
	c.AddSnippet("mtpl_data_st", "$mm = array(")
	c.AddSnippet("mtpl_data_en", ");")
	c.AddSnippet("mtpl_run", "$r = $mmo->add($mm);")
	c.AddSnippet("mtpl_run_check", c.Indent+"if ($r) {"+c.LineBreak+c.Indent+c.Indent+"echo 'Added MailTemplate with ID: ' . $r . PHP_EOL;"+c.LineBreak+c.Indent+"} else {"+c.LineBreak+c.Indent+c.Indent+"echo 'Error adding MailTemplate: ' . $mmo->LAST_ERROR  . PHP_EOL;"+c.LineBreak+c.Indent+"}")

	/* common, group 2 */
	c.AddSnippet("done", "echo 'Done!' . PHP_EOL;")
}

func generateUfTemplate(c *Codegen) string {
	entityFields := map[string]string{}
	var resultStr, snpStr string

	for _, v := range []string{"php", "bxheader", "uf_obj", "uf_data_st"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	entityFields["ENTITY_ID"] = EncloseInSingleQuotes("")
	entityFields["FIELD_NAME"] = EncloseInSingleQuotes("_field_name_")
	entityFields["SORT"] = "500"
	entityFields["XML_ID"] = EncloseInSingleQuotes("")
	entityFields["USER_TYPE_ID"] = EncloseInSingleQuotes("string")
	entityFields["SHOW_FILTER"] = EncloseInSingleQuotes("N")
	entityFields["MULTIPLE"] = EncloseInSingleQuotes("N")
	entityFields["MANDATORY"] = EncloseInSingleQuotes("N")
	entityFields["SHOW_IN_LIST"] = EncloseInSingleQuotes("N")
	entityFields["EDIT_IN_LIST"] = EncloseInSingleQuotes("N")
	entityFields["IS_SEARCHABLE"] = EncloseInSingleQuotes("N")
	entityFields["EDIT_FORM_LABEL"] = "array('ru' => '', 'en' => '')"
	entityFields["LIST_COLUMN_LABEL"] = "array('ru' => '', 'en' => '')"
	entityFields["LIST_FILTER_LABEL"] = "array('ru' => '', 'en' => '')"
	entityFields["ERROR_MESSAGE"] = "array('ru' => '', 'en' => '')"
	entityFields["HELP_MESSAGE"] = "array('ru' => '', 'en' => '')"
	entityFields["SETTINGS"] = "array()"

	for k, v := range entityFields {
		resultStr += fmt.Sprintf(DefaultArrayLine, c.Indent, k, v, c.LineBreak)
	}

	for _, v := range []string{"uf_data_en", "uf_data_run", "uf_data_check", "done"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	return resultStr
}

func generateIbpropTemplate(c *Codegen) string {
	entityFields := map[string]string{}
	var resultStr, snpStr string

	for _, v := range []string{"php", "bxheader", "iblock", "ibp_obj", "iblock_prop_st"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	entityFields["IBLOCK_ID"] = EncloseInSingleQuotes("_0_")
	entityFields["NAME"] = EncloseInSingleQuotes("_name_")
	entityFields["ACTIVE"] = EncloseInSingleQuotes("Y")
	entityFields["SORT"] = "500"
	entityFields["CODE"] = EncloseInSingleQuotes("_code_")
	entityFields["ROW_COUNT"] = "1"
	entityFields["COL_COUNT"] = "30"
	entityFields["XML_ID"] = EncloseInSingleQuotes("")
	entityFields["DEFAULT_VALUE"] = EncloseInSingleQuotes("")
	entityFields["PROPERTY_TYPE"] = EncloseInSingleQuotes("S")
	entityFields["LIST_TYPE"] = EncloseInSingleQuotes("C")
	entityFields["LINK_IBLOCK_ID"] = EncloseInSingleQuotes("0")
	entityFields["MULTIPLE"] = EncloseInSingleQuotes("N")
	entityFields["WITH_DESCRIPTION"] = EncloseInSingleQuotes("N")
	entityFields["SEARCHABLE"] = EncloseInSingleQuotes("N")
	entityFields["FILTRABLE"] = EncloseInSingleQuotes("N")
	entityFields["IS_REQUIRED"] = EncloseInSingleQuotes("N")
	/* Some auto properties */
	entityFields["VERSION"] = "2"
	entityFields["USER_TYPE"] = "false"
	entityFields["USER_TYPE_SETTINGS"] = "false"
	entityFields["HINT"] = EncloseInSingleQuotes("")

	for k, v := range entityFields {
		resultStr += fmt.Sprintf(DefaultArrayLine, c.Indent, k, v, c.LineBreak)
	}

	for _, v := range []string{"iblock_prop_en", "iblock_prop_run", "iblock_prop_check", "done"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	return resultStr
}

func generateMaileventTemplate(c *Codegen) string {
	entityFields := map[string]string{}
	var resultStr, snpStr string

	for _, v := range []string{"php", "bxheader", "mevent_obj", "mevent_data_st"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	entityFields["EVENT_NAME"] = EncloseInSingleQuotes("_event_name_")
	entityFields["LID"] = EncloseInSingleQuotes("ru")
	entityFields["NAME"] = EncloseInSingleQuotes("_name_")
	entityFields["DESCRIPTION"] = EncloseInSingleQuotes("_descr_")
	entityFields["SORT"] = "150"

	for k, v := range entityFields {
		resultStr += fmt.Sprintf(DefaultArrayLine, c.Indent, k, v, c.LineBreak)
	}

	for _, v := range []string{"mevent_data_en", "mevent_run", "mevent_run_check", "mevent_succ"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	for _, v := range []string{"mtpl_obj", "mtpl_data_st", "mtpl_warn"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += c.Indent + snpStr + c.LineBreak
		}
	}

	entityFields = map[string]string{}
	entityFields["EVENT_NAME"] = EncloseInSingleQuotes("_event_name_")
	entityFields["LID"] = EncloseInSingleQuotes("_SID_")
	entityFields["ACTIVE"] = EncloseInSingleQuotes("Y")
	entityFields["EMAIL_FROM"] = EncloseInSingleQuotes("#DEFAULT_EMAIL_FROM#")
	entityFields["EMAIL_TO"] = EncloseInSingleQuotes("#EMAIL_TO#")
	entityFields["SUBJECT"] = EncloseInSingleQuotes("#SUBJECT#")
	entityFields["BODY_TYPE"] = EncloseInSingleQuotes("text")
	entityFields["MESSAGE"] = EncloseInSingleQuotes("Message here with #MACROS#")

	for k, v := range entityFields {
		resultStr += fmt.Sprintf(DefaultArrayLine, c.Indent+c.Indent, k, v, c.LineBreak)
	}

	for _, v := range []string{"mtpl_data_en", "mtpl_run"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += c.Indent + snpStr + c.LineBreak
		}
	}

	for _, v := range []string{"mtpl_run_check", "mevent_run_check_else", "done"} {
		snpStr, _ = c.GetSnippet(v)
		if "" != snpStr {
			resultStr += snpStr + c.LineBreak
		}
	}

	return resultStr
}
