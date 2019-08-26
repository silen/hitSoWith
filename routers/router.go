// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package routers

import (
	"hitSoWith/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//自动路由
	beego.AutoRouter(&controllers.Test{})

}
