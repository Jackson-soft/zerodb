const tools = {
  /**
   * 格式化数字
   * @param s:带格式化数据
   * @param n:保留小数点后几位(默认2位)
   */
  formateNumber(s, n) {
    n = n >= 0 && n <= 20 ? n : 2
    s = `${parseFloat(`${s}`.replace(/[^\d\.-]/g, '')).toFixed(n)}`
    const l = s
      .split('.')[0]
      .split('')
      .reverse()

    const r = s.split('.')[1]
    let t = ''
    for (let i = 0; i < l.length; i++) {
      t += l[i] + ((i + 1) % 3 == 0 && i + 1 != l.length ? ',' : '')
    }
    if (n == 0)
      return t
        .split('')
        .reverse()
        .join('')
    return `${t
      .split('')
      .reverse()
      .join('')}.${r}`
  },
  /**
   * 格式化金额
   * @param s:带格式化数据
   * @param n:保留小数点后几位(默认2位)
   */
  formateMoney(s, n) {
    n = n >= 0 && n <= 20 ? n : 2
    return s.toFixed(n)
  },
  /**
   * 日期格式化
   * @param  date:日期
   * @param  fmt:数据格式
   */

  dateFormate(date, fmt) {
    const o = {
      'M+': date.getMonth() + 1, // 月份
      'd+': date.getDate(), // 日
      'h+': date.getHours(), // 小时
      'm+': date.getMinutes(), // 分
      's+': date.getSeconds(), // 秒
      'q+': Math.floor((date.getMonth() + 3) / 3), // 季度
      S: date.getMilliseconds() // 毫秒
    }
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, `${date.getFullYear()}`.substr(4 - RegExp.$1.length))
    for (const k in o)
      if (new RegExp(`(${k})`).test(fmt))
        fmt = fmt.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k] : `00${o[k]}`.substr(`${o[k]}`.length))
    return fmt
  },
  /**
   * formdate参数转化
   * @param o:需要转化的对象
   * return string: key=value&key=value
   */
  formDate(o) {
    let p = ''
    for (const i in o) {
      p += `${i}=${o[i]}&`
    }
    p = p.substring(0, p.length - 1)
    return p
  },
  /**
   * scrollTo 滚动至某位置
   * @param id || class 需要跳转到的元素位置绑定的id或class
   * demo
   * <div id='t'>...</div>
   * <div class='t'>...</div> 会跳转到文档流中第一个匹配的class
   * scrollTo('t')
   */
  scrollTo(t) {
    const tmp_id = `#${t}`
    const tmp_class = `.${t}`
    let tmp_target = ''

    if (document.querySelector(tmp_id)) {
      tmp_target = document.querySelector(tmp_id)
    } else {
      tmp_target = document.querySelector(tmp_class)
    }

    const targetTop = tmp_target.offsetTop || 0
    window.scrollTo(0, targetTop)
  },
  /**
   * 获取当前页面名称
   */
  getPageName() {
    const path = window.location.pathname
    const pagelist = path.split('/')
    const page = pagelist[pagelist.length - 1].slice(0, -5)
    return page
  },
  /**
   * 获取肉串路由值
   * @param i: 第几个值
   * example: #/first_second_third ====>> 1=>first 2=>second 3=>third
   */
  getRouter(i) {
    const hash = window.location.hash.slice(2)
    const items = hash.split('_')
    if (items.length < 2) {
      return items[0]
    }
    if (i > 0) {
      return items[i - 1]
    }
    return items[0]
  },

  randomString(len, charSet) {
    charSet = charSet || 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let randomString = ''
    for (let i = 0; i < len; i++) {
      const randomPoz = Math.floor(Math.random() * charSet.length)
      randomString += charSet.substring(randomPoz, randomPoz + 1)
    }
    return randomString
  },

  getUrlParam(name) {
    const reg = new RegExp(`(^|&)${name}=([^&]*)(&|$)`) // 构造一个含有目标参数的正则表达式对象
    const r = window.location.search.substr(1).match(reg) // 匹配目标参数
    if (r != null) return unescape(r[2])
    return null // 返回参数值
  }
}

export default tools
