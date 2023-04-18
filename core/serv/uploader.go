package serv

import (
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"strings"
)

type Uploader struct {
}

func NewUploader() *Uploader {
	return &Uploader{}
}

func (my *Uploader) Upload() error {
	return nil
}

func (my *Uploader) Login() error {
	client := resty.New()
	p := url.Values{
		"loginId":          {"15210203617"},
		"password2":        {"cb60bbe75e1cb47755c9237248fa0da3d8b9aece97b49814615b54822dec78bebb99d3915fb3cb2a6ebdf50913bf336443e09b9e65d442bd4824a668290c773b396ebc9af71baa4ee57bfd551d73111cd43d42c3451c97229b4a49bda614b2927fcd1696d8ff5475a1e89f2fbd2f2eab8416cb7e52b794f67517daa0230760e3"},
		"keepLogin":        {"false"},
		"ua":               {"140#oQQDQiTuzzFfvQo2mxXsfpSogBoppq2DxNEGWIrRX0q1yA8xY5A5zJBHTZ0AzMIDKLAUKGU/+qjAcZEBedE8jpAcb9kQLFJWlp1zzqjD6onV1bzx1xwxatJjzzrb22U3lp1xzYqyVDMI2FrrLPc3L6gqzzrVnc3IlpMoeQruRjNIPtrPHpcDk32F85VRIAo+fLvEmg5oLG8arugSOnxfctComAzp/qm5EY4av9pg48AxM6oqK1PGIiELmirXgB93dTsu83VpBP3QGNWkmUJO9Ep3UPIOt3HWFoVSS87v5KJ6YZuHabZ+vdHRmZZnC58uItWXTVeTiqr7xGCZX9lR1MxKe1OxiIE52UlXvirxEZWUZ1IJx4sCUEBzVPMUn7W+y/fZsaWd6shHbBr50kWpuPYcR9j+HRHUBMQn2rSDoBmfBfG78kFrODD0rx/azpwQyTQAhsWAGcr6XMwYh2leHMjK9ORH2kRDIqAY3WCwqMGSLDG18mUwjk4Z919NILRAoGWpbL57U3DlRkuvY/USXL9+4E5Bh4wgjrtpOIKb6LhqyRrs8MtJdWA8cM1gMeGSWhXRx3GPSe4Oc215hNrn/frUSvi6O1gLnM306oRkMBmlzvvhZHhS5H/4gI1dV180oXmTAPeKOM27EiG6ha3w1biHgAyvwm4+90UXTzOSsbaC+YZEZg53Ty5V/uu1VIzhBZyIlTOPhhIkBV7dToNVOOy3wvvin6cD+w+iX5cwhuk1vukVzRFUeYt1Kinp31rnRQLePWOSnzD1jwDg2NmQHe2DoWC8RHIt8Zow7zMni1JRZsbOMZKlnqwpeRx28A4rJvf7Dph4LrrVEjJjgGKUEJFpymi0pLFuX01crWFV2iDoRAW9Xtyp04iSyiHsZeZbJqHbUo6Mn9Ug7fhEEu4OVjZuow8Heqerwg/jciaLvcv7D4Rgs6dlEvOLQfzUT6vsPWBoiOeBmGtuLQQcsWbhNrn2SPWM47hsrNHKLAtU40ucHEBvZQrBbFdhY25CYJ4u4TeuzBnY8lhGpl3eFm+X0M7ldCC9DoovSYnmfOVGDwJuIdRbT24+GAZ5WJgq4SRVPTIDMvKxsUEB5Q0GJIcc7uOQIFVZrFpMJ4It4Tl+d3Ltk2azU3zDOz/9zJQ0TmUAfuaYHRFQDLG+0AUGkJupl5Nk58GNGMRn2PCS85+Yng/0NE/h9m6Qkt3/HqjSozQ="},
		"umidGetStatusVal": {"255"},
		"screenPixel":      {"1920x1080"},
		"navlanguage":      {"zh-CN"},
		"navUserAgent":     {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34"},
		"navPlatform":      {"MacIntel"},
		"appName":          {"taobao"},
		"appEntrance":      {"taobao_h5"},
		"_csrf_token":      {"Z794xMnv5ZMGkHB7KAjCVJ"},
		"umidToken":        {"a51a7f52dc07f1525040dd06616ebc9e18647a28"},
		"hsiz":             {"1c0cd428bdf83d46347f489e23f82baf"},
		"ttid":             {"h5@iframe"},
		"bizParams":        {""},
		"isMobile":         {"true"},
		"lang":             {"zh_CN"},
		"returnUrl":        {"https://h5.m.taobao.com/other/loginend.html?origin=https%3A%2F%2Fmarket.m.taobao.com"},
		"fromSite":         {"0"},
		"umidTag":          {"SERVER"},
		"weiBoMpBridge":    {""},
		"deviceId":         {"kBV2HNQLUwYCAWonEnunBqX1"},
		"pageTraceId":      {"213e203d16818056404456423ead5e"},
		"bx-ua":            {"225!rx6y5ozWooiK 1u7DpDoV dXi lRa/V5co5V0T4Vath1Glxah877Nq5ZJCiNlaF3hnfBQA6PyyyA1hPbFJ8vCEtP1sUAASsWfqgutNTTTOoYtEdjoWVVoJiRDG/ f460Qk0UDMzhfiN2TEVuqHfeSJ3wDGHzfeGeGiCjh/8AfkE30FcGjcI4oL3jDGHzfeG0QESdEbX3NZ dLtpOYB RoL3RDfHkBaa2rjDr485MMIJBoDC1mLsRS9RogvrWxhaLCGok9sNQOf5CbU9EVxf4oJiRDGH fej0Qz0KEM5hf/HAI2oQjsN9Vhii9I5vWbCeAUiOeY QOdt0iI9ql7Omn9hietJMPJ3WIapObAg8CtB5T1pmAAFZSnqL9tN866jrmaj2Xr18f1f5ikSq2grzhqaEVvJG6QuOlWcfyX5OOkH5Fk3SYVJse5MzKtzMxM3rlpuf/XENGIX3oN6Zw8IhVnTUSvXMPhaRvGVOeL5GmdZiizT6B/J55wTh5OZNMYRvvuuYbgPMp1NaDvhqq7K/HnxF5OkxBbVwOpafyyPOuNthdvjurJfFeeiUhN7Uf/woQzaKDM5hfoz6gU4KjcIhoKJkg6VI2iJmN5rso9cLSrf2SQiziwikhoaQfmG7gmLlJJ5XI31cugRxLj3AWFcWLLetTN cU/Pq3182moiua2FTmE5UEy/SF5jFCz8QBlNbY71iG8R/RMl 4Ugg2MjyWFjtaklUxsQVUuvksqeJytlF2tokQH6mwkIc9eolO0LE2oHA39PKeGafeTwzpBBK47J8UbSFhFjErxuZ18m2ZAvOQVfbPLdkURqKa8nPPRa6yMkAM8 MMBi35ur2AE6QskVyErSx38SkgMZo0LD0h8JP3Buefg5F/gtk6fKhYAIRTNzS1hna1DHcI 1PnV5o voFMxJDnEWBMrec lG8O97ka/bxR/ffu8m4mb3gBv8Ma/lGs8/9Z6ZUtWpwIh6OJ I2i/0EFb0rXatvNDZLWY5ux1a4U4XVqLNstbLjVX8cztlSDAkyj3s2GzbTqf 5vernabLijRbzWW9IVmNDBPqDHkpRoM55ufgpT8rnFu3xvtb1T7r1VgGvvRvqvVs5GpD3Pssfk2OpGBZFSeH6Bu/OesgTxN7LN4UltPGzi9k7ZgngVMbDITPcRRKzPc88og3ad/ketbceEvhyolEENFAfXyDaGzR KPBMBM16CaxYGlC7k5SuDdicQvRQ7Np2lRUtFZuB6VOGYSJfl FD2r71hSf V1LOeGqGQrHTQeMrw3NW5INoVIN1QfsoJLaowl8bLnGzyFNGDsWFyKchGCkbBAbLf2nW4MNybBMKrPnflCMb0oRXmiTwnEpeUWST0KC4WcuGE2PhegtKaz3er2o3byqmaBtH3KYJpD4MsJiWDKFChmpd9GQPtjvx8y9OD6Awi284dvWp5NRS0S3d/G2uLJMIZGR2Ure13tYSfaLVe67N2063i1VtY/I//afuiOOZ2 Ez4nCNOo9G/ewJCF3Y2uSy7 TDm0oJ47uWqsodi7RoAcEEsdjGigPyqJGxCMxKyemMlIsaJanftHIUnmkgKPaVqMS1fl5jl6g ArwSzH6Tdswx0XEot1d5RAbh2QS/YATwzchz6X1dxv9VrYKYEyaj6K1MdPtLs4 EdxEKmYB7anMbwopLygXRACb/QX/WptbR2Bdm3stJrHzvNKS KA6/F63pY334Ee7LN5rGW/oKDNkeQbRY3CS CztKZYgnHaQmdoRdBCX0o9cx2OlwSxrPIH2KwuFQ02Jzfw3XzYDO 5jUTtwQsfktmIMCe59CFGkPvfvu2s3ePKUhNa5Kl6bCWU5yM0XWF65H 2GLBIaPjumDGg7qKttaQZN6nCcNiKvG67vuGxdQgVFJZ5OtwuBtkhhVnwB6fL6vsLqsiI6t8nX4BhXK3M2PIP09X Jj4n8B6rgfO2IzNCdHf/aIB0 jKdo5I0JDQkI1eFElOuMjr1TyZrjB 7DBl19hZX6wHjJUw9/GPlSeXt42MlcvF3UpgW2rwGgSQDZy9/5ycWpoLF3BFsJUgl23Ene4ck90 vARYb="},
		"bx-umidtoken":     {"GF74D17E75D84F7D48F5E4E68CE6C59ADC733FDBB7D1FDE968F"},
	}
	c := []*http.Cookie{
		{Name: "co", Value: "2"},
		{Name: "XSRF-TOKEN", Value: "60f9e3df-151b-418f-9427-aaa03fb3aa6a"},
		{Name: "_uab_collina", Value: "168128254021668805695683"},
		{Name: "_bl_uid", Value: "L0lnkgIqd1ecj87wyq7XwXvip3ya"},
		{Name: "arms_uid", Value: "13259823-388f-4d75-bf56-4123d1c287ce"},
		{Name: "cna", Value: "kBV2HNQLUwYCAWonEnunBqX1"},
		{Name: "xlly_s", Value: "1"},
		{Name: "_samesite_flag_", Value: "true"},
		{Name: "cookie2", Value: "1e2f8db89b5d7c385193b2aba5f0c3c2"},
		{Name: "t", Value: "e38777c23c69bfa31c78601437ed7526"},
		{Name: "_tb_token_", Value: "fe751874d3fa"},
		{Name: "x5sec", Value: "7b22617365727665723b32223a226633653737616134643531616663623033626664646239656165376639346632434b32722b614547454f2f79784b6a4f3938476a69514561445445314d6a45774d6a417a4e6a45334f7a45777449547030415a4141773d3d227d"},
		{Name: "isg", Value: "BPz8CGPrtLGHSoB2qy3urhdJzZyu9aAfZHKL7tZ8Q-fKoZsr_AXOrhFXgcnZ6dh3"},
		{Name: "l", Value: "fBTEQ83nNxFc2gWUBO5Churza77TqIOffsPzaNbMiIEGa63PTFiTZNCs8lVe7dtjgTfqyexrSKexYdFXJba3Wxt2kUkmCzwjq19M-2R7Zb1.."},
		{Name: "tfstk", Value: "dh-kA12BhsNsth--JgISb2gRva0AN7sC_BEd9MCEus5fv_KJyeYDMsEKFg3CxBfDO_QR4Bd3tQdA8BCpe8TX6drJ24sLxgsCYfhtXc3WFMsUQyeZkTb73g-Uxcn9FLZFOJkEXXT4DkS9Hkz5wAKxENCghxH0-EPsI6pPilvXcTzd6c18YfKAEKfNQta47sTUAtkpmyaCzt6c68_E8SC.."},
	}

	uri, _ := url.Parse("https://login.m.taobao.com/newlogin/login.do?appName=taobao&fromSite=0&_bx-v=2.2.3")
	uri.RawQuery = p.Encode()
	res, err := client.R().SetCookies(c).
		SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34").
		SetHeader("referer", "https://login.m.taobao.com/login.htm?ttid=h5%40iframe&redirectURL=https%3A%2F%2Fh5.m.taobao.com%2Fother%2Floginend.html%3Forigin%3Dhttps%253A%252F%252Fmarket.m.taobao.com").
		Post(uri.String())
	if err != nil {
		return err
	}
	values := res.Header().Values("Set-Cookie")
	for _, v := range values {
		pair := strings.Split(strings.SplitN(v, ";", 2)[0], "=")
		println(strings.Trim(pair[0], ""))
		println(strings.Trim(pair[1], ""))
	}
	println(res.String())
	return nil
}

func (my *Uploader) signature(params, cookies map[string]string, retryTime int) error {
	client := resty.New()
	p := url.Values{}
	var c []*http.Cookie
	for k, v := range params {
		p.Add(k, v)
	}
	for k, v := range cookies {
		c = append(c, &http.Cookie{
			Name: k, Value: v,
		})
	}
	uri, _ := url.Parse("https://h5api.m.taobao.com/h5/mtop.taobao.livexadmin.datainf.oss.signature/1.0/")
	uri.RawQuery = p.Encode()
	res, err := client.R().SetCookies(c).
		SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.34").
		SetHeader("referer", "https://market.m.taobao.com/").Get(uri.String())
	if err != nil {
		return err
	}

	values := res.Header().Values("Set-Cookie")
	for _, v := range values {
		pair := strings.Split(strings.SplitN(v, ";", 2)[0], "=")
		cookies[strings.Trim(pair[0], "")] = strings.Trim(pair[1], "")
	}
	ref := gjson.GetBytes(res.Body(), "ret.0").String()

	if retryTime < 10 && ("FAIL_SYS_SESSION_EXPIRED::Session过期" == ref || "FAIL_SYS_TOKEN_EMPTY::令牌为空" == ref) {
		err := my.signature(params, cookies, retryTime+1)
		if err != nil {
			return err
		}
	}
	println(res.String())

	//in, err := util.ReadExcel("/Users/Chaly/Documents/Workspace/yugong/docs/header.xlsx", "params", "cookies")
	//if err != nil {
	//	return err
	//}
	//params := url.Values{}
	//var cookies []*http.Cookie
	//keys := map[string]bool{
	//	"appKey": true, "t": true, "sign": true, "data": true,
	//	//"cookie2": true, "sg": true, "_m_h5_tk": true, "_m_h5_tk_enc": true,
	//}
	//for _, p := range in["params"] {
	//	if _, ok := keys[p[0]]; ok {
	//		params.Add(p[0], p[1])
	//		continue
	//	}
	//}
	//for _, c := range in["cookies"] {
	//	if _, ok := keys[c[0]]; ok {
	//		cookies = append(cookies, &http.Cookie{
	//			Name:  strings.Trim(c[0], ""),
	//			Value: strings.Trim(c[1], ""),
	//		})
	//	}
	//}

	//过期会set-cookie2
	//https://h5api.m.taobao.com/h5/mtop.taobao.livexadmin.datainf.oss.signature/1.0/
	//?jsv=2.4.16&appKey=12574478&t=1681797935061&sign=8085c963ed7c3906fc1b6cef217322e9
	//&api=mtop.taobao.livexadmin.datainf.oss.signature&v=1.0&type=jsonp&dataType=jsonp&callback=mtopjsonp1&data=%7B%7D
	//cookie2=1821f17fa288732b6bbb351595b60a16;Path=/;Domain=.taobao.com;Max-Age=-1;HttpOnly
	//{"api":"mtop.taobao.livexadmin.datainf.oss.signature","data":{},"ret":["FAIL_SYS_SESSION_EXPIRED::Session过期"],"v":"1.0"}

	//cookie2 = strings.SplitN(cookie2, ";", 2)[0]

	//if "FAIL_SYS_SESSION_EXPIRED::Session过期" == ref {
	//	err := my.signature(cookie2)
	//	if err != nil {
	//		return err
	//	}
	//}

	//令牌为空会set-_m_h5_tk&set-_m_h5_tk_enc
	//{"api":"mtop.taobao.livexadmin.datainf.oss.signature","data":{},"ret":["FAIL_SYS_TOKEN_EMPTY::令牌为空"],"v":"1.0"}

	//if "FAIL_SYS_TOKEN_EMPTY::令牌为空" == gjson.GetBytes(res.Body(), "ret.0").String() {
	//	_m_h5_tk := res.RawResponse.Header.Get("Set-Cookie")
	//	println(_m_h5_tk)
	//	//err := my.signature(_m_h5_tk)
	//	//if err != nil {
	//	//	return err
	//	//}
	//}

	//println(res.RawResponse.Header.Get("Set-Cookie"))
	//for s := range res.RawResponse.Header {
	//	println(s, res.RawResponse.Header.Get(s))
	//}
	//亮2b
	//res.RawResponse.Header.Get("Set-Cookie")
	return nil
}

//https://login.m.taobao.com/newlogin/login.do?appName=taobao&fromSite=0&_bx-v=2.2.3
//https://market.m.taobao.com/app/tb-zhibo-app/mcn-center/daren.html#/

//https://h5api.m.taobao.com/h5/mtop.taobao.livexadmin.datainf.oss.signature/1.0/?jsv=2.4.16&appKey=12574478&t=1681789811698&sign=74dfd668f59da6dcc3415ae77e85f5e5&api=mtop.taobao.livexadmin.datainf.oss.signature&v=1.0&type=jsonp&dataType=jsonp&callback=mtopjsonp1&data=%7B%7D
