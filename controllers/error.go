package controllers

func HandleError(c *TableauController, status uint, errormsg string) {
	c.Ctx.Output.SetStatus(int(status))
	c.Data["json"] = map[string]string{"error": errormsg}
	c.ServeJSON()
}
