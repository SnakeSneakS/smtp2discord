package main_test

import (
	"fmt"
	"net/smtp"
	"strings"
	"testing"
	"time"

	. "github.com/snakesneaks/smtp2discord"
)

func TestMain(t *testing.T) {
	server := NewSmtp2DiscordServer()
	// サーバーの起動
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !IsClosedNetworkError(err) {
				t.Errorf("Failed to start server: %v", err)
			}
		}
	}()

	// サーバーが起動するまで待機
	time.Sleep(100 * time.Millisecond)

	const sender = "sender@example.com"
	const recipient = "recipient@example.com"
	//example message like Nextcloud
	const message = "Subject: SUBJECT\r\nTo: TO \u003cTO@example.com\u003e\r\nFrom: From \u003cfrom@example.com\u003e\r\nMessage-ID: \u003cb7dd792e4973022285bffc23fcd1fcfa@example.com\u003e\r\nMIME-Version: 1.0\r\nDate: Thu, 26 Dec 2024 09:49:01 +0000\r\nContent-Type: multipart/alternative; boundary=uI7wL8bs\r\n\r\n--uI7wL8bs\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\ntest=E3=81=95=E3=82=93=E3=80=81=E6=88=90=E5=8A=9F=E3=81=A7=\r\n=E3=81=99=EF=BC=81\r\n\r\n=E3=81=93=E3=81=AE=E3=83=A1=E3=83=BC=E3=83=AB=\r\n=E3=81=8C=E8=A6=8B=E3=81=88=E3=81=A6=E3=81=84=E3=82=8B=E3=81=A8=E3=81=84=\r\n=E3=81=86=E3=81=93=E3=81=A8=E3=81=AF=E3=80=81=E3=83=A1=E3=83=BC=E3=83=AB=\r\n=E8=A8=AD=E5=AE=9A=E3=81=8C=E3=81=86=E3=81=BE=E3=81=8F=E3=81=84=E3=81=A3=\r\n=E3=81=A6=E3=81=84=E3=82=8B=E3=81=A8=E3=81=84=E3=81=86=E3=81=93=E3=81=A8=\r\n=E3=81=A7=E3=81=99=E3=80=82\r\n\r\n\r\n--=20\r\nNextcloud - =E3=81=99=E3=81=B9=\r\n=E3=81=A6=E3=81=AE=E3=83=87=E3=83=BC=E3=82=BF=E3=82=92=E5=AE=89=E5=85=A8=\r\n=E3=81=AB=E4=BF=9D=E7=AE=A1=E3=81=97=E3=81=BE=E3=81=99\r\n=E3=81=93=E3=82=\r\n=8C=E3=81=AF=E8=87=AA=E5=8B=95=E7=9A=84=E3=81=AB=E7=94=9F=E6=88=90=E3=81=\r\n=95=E3=82=8C=E3=81=9F=E3=83=A1=E3=83=BC=E3=83=AB=E3=81=A7=E3=81=99=E3=80=\r\n=82=E8=BF=94=E4=BF=A1=E3=81=97=E3=81=AA=E3=81=84=E3=81=A7=E3=81=8F=E3=81=\r\n=A0=E3=81=95=E3=81=84=E3=80=82\r\n--uI7wL8bs\r\nContent-Type: text/html; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n\u003c!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.=\r\nw3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"\u003e\r\n\u003chtml xmlns=3D\"http://www.=\r\nw3.org/1999/xhtml\" lang=3D\"en\" xml:lang=3D\"en\" style=3D\"-webkit-font-smooth=\r\ning:antialiased;background:#fff!important\"\u003e\r\n\u003chead\u003e\r\n=09\u003cmeta http-equiv=\r\n=3D\"Content-Type\" content=3D\"text/html; charset=3Dutf-8\"\u003e\r\n=09\u003cmeta name=\r\n=3D\"viewport\" content=3D\"width=3Ddevice-width\"\u003e\r\n=09\u003ctitle\u003e\u003c/title\u003e\r\n=09\u003c=\r\nstyle type=3D\"text/css\"\u003e@media only screen{html{min-height:100%;background:=\r\n#fff}}@media only screen and (max-width:610px){table.body img{width:auto;he=\r\night:auto}table.body center{min-width:0!important}table.body .container{wid=\r\nth:95%!important}table.body .columns{height:auto!important;-moz-box-sizing:=\r\nborder-box;-webkit-box-sizing:border-box;box-sizing:border-box;padding-left=\r\n:30px!important;padding-right:30px!important}th.small-12{display:inline-blo=\r\nck!important;width:100%!important}table.menu{width:100%!important}table.men=\r\nu td,table.menu th{width:auto!important;display:inline-block!important}tabl=\r\ne.menu.vertical td,table.menu.vertical th{display:block!important}table.men=\r\nu[align=3Dcenter]{width:auto!important}}\u003c/style\u003e\r\n\u003c/head\u003e\r\n\u003cbody style=3D=\r\n\"-moz-box-sizing:border-box;-ms-text-size-adjust:100%;-webkit-box-sizing:bo=\r\nrder-box;-webkit-font-smoothing:antialiased;-webkit-text-size-adjust:100%;m=\r\nargin:0;background:#fff!important;box-sizing:border-box;color:#0a0a0a;font-=\r\nfamily:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubunt=\r\nu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:16px;font-weight:40=\r\n0;line-height:1.3;margin:0;min-width:100%;padding:0;text-align:left;width:1=\r\n00%!important\"\u003e\r\n=09\u003cspan class=3D\"preheader\" style=3D\"color:#F5F5F5;displ=\r\nay:none!important;font-size:1px;line-height:1px;max-height:0;max-width:0;ms=\r\no-hide:all!important;opacity:0;overflow:hidden;visibility:hidden\"\u003e\r\n=09\u003c/s=\r\npan\u003e\r\n=09\u003ctable class=3D\"body\" style=3D\"-webkit-font-smoothing:antialiased=\r\n;margin:0;background:#fff;border-collapse:collapse;border-spacing:0;color:#=\r\n0a0a0a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxyge=\r\nn-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:16px;fo=\r\nnt-weight:400;line-height:1.3;margin:0;padding:0;text-align:left;vertical-a=\r\nlign:top;width:100%\"\u003e\r\n=09=09\u003ctr style=3D\"padding:0;text-align:left;vertic=\r\nal-align:top\"\u003e\r\n=09=09=09\u003ctd class=3D\"center\" align=3D\"center\" valign=3D\"t=\r\nop\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Margin:0;border-collaps=\r\ne:collapse!important;color:#0a0a0a;font-family:-apple-system,BlinkMacSystem=\r\nFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,=\r\nsans-serif;font-size:16px;font-weight:400;hyphens:auto;line-height:1.3;marg=\r\nin:0;padding:0;text-align:left;vertical-align:top;word-wrap:break-word\"\u003e\r\n=\r\n=09=09=09=09\u003ccenter data-parsed=3D\"\" style=3D\"min-width:580px;width:100%\"\u003e\u003c=\r\ntable align=3D\"center\" class=3D\"wrapper header float-center\" style=3D\"Margi=\r\nn:0 auto;background:#fff;border-collapse:collapse;border-spacing:0;float:no=\r\nne;margin:0 auto;padding:0;text-align:center;vertical-align:top;width:100%\"=\r\n\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=09=\r\n\u003ctd class=3D\"wrapper-inner\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto=\r\n;Margin:0;border-collapse:collapse!important;color:#0a0a0a;font-family:-app=\r\nle-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell=\r\n,'Helvetica Neue',Arial,sans-serif;font-size:16px;font-weight:400;hyphens:a=\r\nuto;line-height:1.3;margin:0;padding:20px;text-align:left;vertical-align:to=\r\np;word-wrap:break-word\"\u003e\r\n=09=09=09\u003ctable align=3D\"center\" class=3D\"contai=\r\nner\" style=3D\"Margin:0 auto;background:0 0;border-collapse:collapse;border-=\r\nspacing:0;margin:0 auto;padding:0;text-align:inherit;vertical-align:top;wid=\r\nth:150px\"\u003e\r\n=09=09=09=09\u003ctbody\u003e\r\n=09=09=09=09\u003ctr style=3D\"padding:0;text-=\r\nalign:left;vertical-align:top\"\u003e\r\n=09=09=09=09=09\u003ctd style=3D\"-moz-hyphens:=\r\nauto;-webkit-hyphens:auto;Margin:0;border-collapse:collapse!important;color=\r\n:#0a0a0a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxy=\r\ngen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:16px;=\r\nfont-weight:400;hyphens:auto;line-height:1.3;margin:0;padding:0;text-align:=\r\nleft;vertical-align:top;word-wrap:break-word\"\u003e\r\n=09=09=09=09=09=09\u003ctable c=\r\nlass=3D\"row collapse\" style=3D\"border-collapse:collapse;border-spacing:0;di=\r\nsplay:table;padding:0;position:relative;text-align:left;vertical-align:top;=\r\nwidth:100%\"\u003e\r\n=09=09=09=09=09=09=09\u003ctbody\u003e\r\n=09=09=09=09=09=09=09\u003ctr styl=\r\ne=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=09=09=09=09=09=09=\r\n=09\u003ccenter data-parsed=3D\"\" style=3D\"background-color:#00679e;min-width:175=\r\npx;max-height:175px; padding:35px 0px;border-radius:200px\"\u003e\r\n=09=09=09=09=\r\n=09=09=09=09=09\u003cimg class=3D\"logo float-center\" src=3D\"https://nextcloud.my=\r\n.snakesneaks.xyz/core/img/logo/logo.png?v=3D0\" alt=3D\"Nextcloud\" align=3D\"c=\r\nenter\" style=3D\"-ms-interpolation-mode:bicubic;clear:both;display:block;flo=\r\nat:none;margin:0 auto;outline:0;text-align:center;text-decoration:none;max-=\r\nheight:105px;max-width:105px;width:auto;height:auto\" width=3D\"105\" height=\r\n=3D\"50\"\u003e\r\n=09=09=09=09=09=09=09=09\u003c/center\u003e\r\n=09=09=09=09=09=09=09\u003c/tr\u003e=\r\n\r\n=09=09=09=09=09=09=09\u003c/tbody\u003e\r\n=09=09=09=09=09=09\u003c/table\u003e\r\n=09=09=09=\r\n=09=09\u003c/td\u003e\r\n=09=09=09=09\u003c/tr\u003e\r\n=09=09=09=09\u003c/tbody\u003e\r\n=09=09=09\u003c/table\u003e=\r\n\r\n=09=09\u003c/td\u003e\r\n=09\u003c/tr\u003e\r\n\u003c/table\u003e\r\n\u003ctable class=3D\"spacer float-center\"=\r\n style=3D\"Margin:0 auto;border-collapse:collapse;border-spacing:0;float:non=\r\ne;margin:0 auto;padding:0;text-align:center;vertical-align:top;width:100%\"\u003e=\r\n\r\n=09\u003ctbody\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align:left;vertical-align:to=\r\np\"\u003e\r\n=09=09\u003ctd height=3D\"40px\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:=\r\nauto;Margin:0;border-collapse:collapse!important;color:#0a0a0a;font-size:80=\r\npx;font-weight:400;hyphens:auto;line-height:80px;margin:0;mso-line-height-r=\r\nule:exactly;padding:0;text-align:left;vertical-align:top;word-wrap:break-wo=\r\nrd\"\u003e\u0026#xA0;\u003c/td\u003e\r\n=09\u003c/tr\u003e\r\n=09\u003c/tbody\u003e\r\n\u003c/table\u003e\u003ctable align=3D\"center\" =\r\nclass=3D\"container main-heading float-center\" style=3D\"Margin:0 auto;backgr=\r\nound:0 0!important;border-collapse:collapse;border-spacing:0;float:none;mar=\r\ngin:0 auto;padding:0;text-align:center;vertical-align:top;width:580px\"\u003e\r\n=\r\n=09\u003ctbody\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align:left;vertical-align:top\"\u003e=\r\n\r\n=09=09\u003ctd style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Margin:0;borde=\r\nr-collapse:collapse!important;color:#0a0a0a;font-family:-apple-system,Blink=\r\nMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neu=\r\ne',Arial,sans-serif;font-size:16px;font-weight:400;hyphens:auto;line-height=\r\n:1.3;margin:0;padding:0;text-align:left;vertical-align:top;word-wrap:break-=\r\nword\"\u003e\r\n=09=09=09\u003ch1 class=3D\"text-center\" style=3D\"Margin:0;Margin-bottom=\r\n:10px;color:inherit;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI'=\r\n,Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font=\r\n-size:24px;font-weight:400;line-height:1.3;margin:0;padding:0;text-align:ce=\r\nnter;word-wrap:normal\"\u003eadmin-N9X0y8PrX6rR=E3=81=95=E3=82=93=E3=80=81=\r\n=E6=88=90=E5=8A=9F=E3=81=A7=E3=81=99=EF=BC=81\u003c/h1\u003e\r\n=09=09\u003c/td\u003e\r\n=09\u003c/tr\u003e=\r\n\r\n=09\u003c/tbody\u003e\r\n\u003c/table\u003e\r\n\u003ctable class=3D\"spacer float-center\" style=3D\"M=\r\nargin:0 auto;border-collapse:collapse;border-spacing:0;float:none;margin:0 =\r\nauto;padding:0;text-align:center;vertical-align:top;width:100%\"\u003e\r\n=09\u003ctbod=\r\ny\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=\r\n=09\u003ctd height=3D\"36px\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Marg=\r\nin:0;border-collapse:collapse!important;color:#0a0a0a;font-size:40px;font-w=\r\neight:400;hyphens:auto;line-height:36px;margin:0;mso-line-height-rule:exact=\r\nly;padding:0;text-align:left;vertical-align:top;word-wrap:break-word\"\u003e\u0026#xA0=\r\n;\u003c/td\u003e\r\n=09\u003c/tr\u003e\r\n=09\u003c/tbody\u003e\r\n\u003c/table\u003e\u003ctable align=3D\"center\" class=3D\"=\r\nwrapper content float-center\" style=3D\"Margin:0 auto;border-collapse:collap=\r\nse;border-spacing:0;float:none;margin:0 auto;padding:0;text-align:center;ve=\r\nrtical-align:top;width:100%\"\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align:left;v=\r\nertical-align:top\"\u003e\r\n=09=09\u003ctd class=3D\"wrapper-inner\" style=3D\"-moz-hyphe=\r\nns:auto;-webkit-hyphens:auto;Margin:0;border-collapse:collapse!important;co=\r\nlor:#0a0a0a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,=\r\nOxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:16=\r\npx;font-weight:400;hyphens:auto;line-height:1.3;margin:0;padding:0;text-ali=\r\ngn:left;vertical-align:top;word-wrap:break-word\"\u003e\r\n=09=09=09\u003ctable align=\r\n=3D\"center\" class=3D\"container\" style=3D\"Margin:0 auto;background:#fff;bord=\r\ner-collapse:collapse;border-spacing:0;margin:0 auto;padding:0;text-align:in=\r\nherit;vertical-align:top;width:580px\"\u003e\r\n=09=09=09=09\u003ctbody\u003e\r\n=09=09=09=09=\r\n\u003ctr style=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=09=09=09=\r\n=09\u003ctd style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Margin:0;border-coll=\r\napse:collapse!important;color:#0a0a0a;font-family:-apple-system,BlinkMacSys=\r\ntemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Ari=\r\nal,sans-serif;font-size:16px;font-weight:400;hyphens:auto;line-height:1.3;m=\r\nargin:0;padding:0;text-align:left;vertical-align:top;word-wrap:break-word\"\u003e=\r\n\u003ctable class=3D\"row description\" style=3D\"border-collapse:collapse;border-s=\r\npacing:0;display:table;padding:0;position:relative;text-align:left;vertical=\r\n-align:top;width:100%\"\u003e\r\n=09\u003ctbody\u003e\r\n=09\u003ctr style=3D\"padding:0;text-align=\r\n:left;vertical-align:top\"\u003e\r\n=09=09\u003cth class=3D\"small-12 large-12 columns f=\r\nirst last\" style=3D\"Margin:0 auto;color:#0a0a0a;font-family:-apple-system,B=\r\nlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica=\r\n Neue',Arial,sans-serif;font-size:16px;font-weight:400;line-height:1.3;marg=\r\nin:0 auto;padding:0;padding-bottom:30px;padding-left:30px;padding-right:30p=\r\nx;text-align:left;width:550px\"\u003e\r\n=09=09=09\u003ctable style=3D\"border-collapse:=\r\ncollapse;border-spacing:0;padding:0;text-align:left;vertical-align:top;widt=\r\nh:100%\"\u003e\r\n=09=09=09=09\u003ctr style=3D\"padding:0;text-align:left;vertical-alig=\r\nn:top\"\u003e\r\n=09=09=09=09=09\u003cth style=3D\"Margin:0;color:#0a0a0a;font-family:-a=\r\npple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantare=\r\nll,'Helvetica Neue',Arial,sans-serif;font-size:16px;font-weight:400;line-he=\r\night:1.3;margin:0;padding:0;text-align:left\"\u003e\r\n=09=09=09=09=09=09\u003cp style=\r\n=3D\"Margin:0;Margin-bottom:10px;color:#777;font-family:-apple-system,BlinkM=\r\nacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue=\r\n',Arial,sans-serif;font-size:16px;font-weight:400;line-height:1.3;margin:0;=\r\nmargin-bottom:10px;padding:0;text-align:center\"\u003e=E3=81=93=E3=81=AE=E3=83=\r\n=A1=E3=83=BC=E3=83=AB=E3=81=8C=E8=A6=8B=E3=81=88=E3=81=A6=E3=81=84=E3=82=\r\n=8B=E3=81=A8=E3=81=84=E3=81=86=E3=81=93=E3=81=A8=E3=81=AF=E3=80=81=E3=83=\r\n=A1=E3=83=BC=E3=83=AB=E8=A8=AD=E5=AE=9A=E3=81=8C=E3=81=86=E3=81=BE=E3=81=\r\n=8F=E3=81=84=E3=81=A3=E3=81=A6=E3=81=84=E3=82=8B=E3=81=A8=E3=81=84=E3=81=\r\n=86=E3=81=93=E3=81=A8=E3=81=A7=E3=81=99=E3=80=82\u003c/p\u003e\r\n=09=09=09=09=09\u003c/th\u003e=\r\n\r\n=09=09=09=09=09\u003cth class=3D\"expander\" style=3D\"Margin:0;color:#0a0a0a;fo=\r\nnt-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ub=\r\nuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:16px;font-weight=\r\n:400;line-height:1.3;margin:0;padding:0!important;text-align:left;visibilit=\r\ny:hidden;width:0\"\u003e\u003c/th\u003e\r\n=09=09=09=09\u003c/tr\u003e\r\n=09=09=09\u003c/table\u003e\r\n=09=09\u003c/t=\r\nh\u003e\r\n=09\u003c/tr\u003e\r\n=09\u003c/tbody\u003e\r\n\u003c/table\u003e\r\n=09=09=09=09=09\u003c/td\u003e\r\n=09=09=09=\r\n=09\u003c/tr\u003e\r\n=09=09=09=09\u003c/tbody\u003e\r\n=09=09=09\u003c/table\u003e\r\n=09=09\u003c/td\u003e\r\n=09\u003c/tr=\r\n\u003e\r\n\u003c/table\u003e\u003ctable class=3D\"spacer float-center\" style=3D\"Margin:0 auto;bor=\r\nder-collapse:collapse;border-spacing:0;float:none;margin:0 auto;padding:0;t=\r\next-align:center;vertical-align:top;width:100%\"\u003e\r\n=09\u003ctbody\u003e\r\n=09\u003ctr styl=\r\ne=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=09\u003ctd height=3D\"6=\r\n0px\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Margin:0;border-collap=\r\nse:collapse!important;color:#0a0a0a;font-family:-apple-system,BlinkMacSyste=\r\nmFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial=\r\n,sans-serif;font-size:60px;font-weight:400;hyphens:auto;line-height:60px;ma=\r\nrgin:0;mso-line-height-rule:exactly;padding:0;text-align:left;vertical-alig=\r\nn:top;word-wrap:break-word\"\u003e\u0026#xA0;\u003c/td\u003e\r\n=09\u003c/tr\u003e\r\n=09\u003c/tbody\u003e\r\n\u003c/table\u003e=\r\n\r\n\u003ctable align=3D\"center\" class=3D\"wrapper footer float-center\" style=3D\"M=\r\nargin:0 auto;border-collapse:collapse;border-spacing:0;float:none;margin:0 =\r\nauto;padding:0;text-align:center;vertical-align:top;width:100%\"\u003e\r\n=09\u003ctr s=\r\ntyle=3D\"padding:0;text-align:left;vertical-align:top\"\u003e\r\n=09=09\u003ctd class=3D=\r\n\"wrapper-inner\" style=3D\"-moz-hyphens:auto;-webkit-hyphens:auto;Margin:0;bo=\r\nrder-collapse:collapse!important;color:#0a0a0a;font-family:-apple-system,Bl=\r\ninkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'Helvetica =\r\nNeue',Arial,sans-serif;font-size:16px;font-weight:400;hyphens:auto;line-hei=\r\nght:1.3;margin:0;padding:0;text-align:left;vertical-align:top;word-wrap:bre=\r\nak-word\"\u003e\r\n=09=09=09\u003ccenter data-parsed=3D\"\" style=3D\"min-width:580px;widt=\r\nh:100%\"\u003e\r\n=09=09=09=09\u003ctable class=3D\"spacer float-center\" style=3D\"Margin=\r\n:0 auto;border-collapse:collapse;border-spacing:0;float:none;margin:0 auto;=\r\npadding:0;text-align:center;vertical-align:top;width:100%\"\u003e\r\n=09=09=09=09=\r\n=09\u003ctbody\u003e\r\n=09=09=09=09=09\u003ctr style=3D\"padding:0;text-align:left;vertical=\r\n-align:top\"\u003e\r\n=09=09=09=09=09=09\u003ctd height=3D\"15px\" style=3D\"-moz-hyphens:=\r\nauto;-webkit-hyphens:auto;Margin:0;border-collapse:collapse!important;color=\r\n:#0a0a0a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxy=\r\ngen-Sans,Ubuntu,Cantarell,'Helvetica Neue',Arial,sans-serif;font-size:15px;=\r\nfont-weight:400;hyphens:auto;line-height:15px;margin:0;mso-line-height-rule=\r\n:exactly;padding:0;text-align:left;vertical-align:top;word-wrap:break-word\"=\r\n\u003e\u0026#xA0;\u003c/td\u003e\r\n=09=09=09=09=09\u003c/tr\u003e\r\n=09=09=09=09=09\u003c/tbody\u003e\r\n=09=09=09=\r\n=09\u003c/table\u003e\r\n=09=09=09=09\u003cp class=3D\"text-center float-center\" align=3D\"ce=\r\nnter\" style=3D\"Margin:0;Margin-bottom:10px;color:#C8C8C8;font-family:-apple=\r\n-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen-Sans,Ubuntu,Cantarell,'=\r\nHelvetica Neue',Arial,sans-serif;font-size:12px;font-weight:400;line-height=\r\n:16px;margin:0;margin-bottom:10px;padding:0;text-align:center\"\u003eNextcloud - =\r\n=E3=81=99=E3=81=B9=E3=81=A6=E3=81=AE=E3=83=87=E3=83=BC=E3=82=BF=E3=82=92=\r\n=E5=AE=89=E5=85=A8=E3=81=AB=E4=BF=9D=E7=AE=A1=E3=81=97=E3=81=BE=E3=81=99\u003cbr=\r\n\u003e=E3=81=93=E3=82=8C=E3=81=AF=E8=87=AA=E5=8B=95=E7=9A=84=E3=81=AB=E7=94=\r\n=9F=E6=88=90=E3=81=95=E3=82=8C=E3=81=9F=E3=83=A1=E3=83=BC=E3=83=AB=E3=81=\r\n=A7=E3=81=99=E3=80=82=E8=BF=94=E4=BF=A1=E3=81=97=E3=81=AA=E3=81=84=E3=81=\r\n=A7=E3=81=8F=E3=81=A0=E3=81=95=E3=81=84=E3=80=82\u003c/p\u003e\r\n=09=09=09\u003c/center\u003e=\r\n\r\n=09=09\u003c/td\u003e\r\n=09\u003c/tr\u003e\r\n\u003c/table\u003e=09=09=09=09=09\u003c/center\u003e\r\n=09=09=09=09=\r\n\u003c/td\u003e\r\n=09=09=09\u003c/tr\u003e\r\n=09=09\u003c/table\u003e\r\n=09=09\u003c!-- prevent Gmail on iOS f=\r\nont size manipulation --\u003e\r\n=09=09\u003cdiv style=3D\"display:none;white-space:no=\r\nwrap;font:15px courier;line-height:0\"\u003e\u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026n=\r\nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; =\r\n\u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp; \u0026nbsp=\r\n; \u0026nbsp; \u0026nbsp; \u0026nbsp;\u003c/div\u003e\r\n=09\u003c/body\u003e\r\n\u003c/html\u003e\r\n\r\n--uI7wL8bs--\r\n\r\n"
	serverURL := fmt.Sprintf("%s:%s", Cfg.Server.Domain, strings.Split(Cfg.Server.Addr, ":")[1])
	auth := smtp.PlainAuth("", Cfg.Auth.Username, Cfg.Auth.Password, "localhost")
	err := smtp.SendMail(serverURL, auth, sender, []string{recipient}, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}
