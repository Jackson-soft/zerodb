export class loading {
  /**
   * @param data:{
              * loadingId:'loading',//div的id
              * background : 'rgba(255,255,255,.5)',//loading背景颜色，（可选）
              * loadingStyle : {
                    * type: 'default',//default（默认值，默认显示css3的加载动画），image设置自定义背景图片
                    * width: '10%',//背景图片大小
                    * imgUrl: 'https://assets.2dfire.com/frontend/56ef0094711c12fafd86e0a2ede914db.gif',
                * },
                * minShowTime: 3000,//loading动画最小显示的时长
     * */
  constructor (data) {
    // console.log(data)
    const defaultConfig = {
      loadingId: 'loading',
      showType: 'default', // 默认全局加载，局部加载为only
      loadingStyle: {
        type: 'default', // default默认值，
        width: '8em',
        imgUrl: 'https://assets.2dfire.com/frontend/56ef0094711c12fafd86e0a2ede914db.gif',
      },
      background: 'rgba(255,255,255,.5)',
      minShowTime: 1000
    }
    this.config = Object.assign(this, defaultConfig, data)

    this.config.loadDOM = document.getElementById(this.config.loadingId)
    this.num = this.config.loadDOM ? Number(this.config.loadDOM.getAttribute('data-num')) : 0

    this.loadingSetTimeout = {}// 用于放置当前loading是否显示，默认loading过3秒之后才能结束

    this._loadingInit()
  }

  /**
   * loading初始化
   * */
  _loadingInit () {
    if (this.config.loadDOM) {
      return false
    }
    this._addDom(this.config)
    this._loadCssCode(this._getDefaultStyle())
    this._addStyle(this.config)
  }

  /**
   * 添加html节点
   * */
  _addDom (config, callback) {
    const div = document.createElement('div')
    div.id = config.loadingId
    let html = `<div class="m-loading-wrapper">`
    html += this._getTypeDom(config.loadingStyle.type)
    html += `</div>`
    div.innerHTML = html

    /** 添加状态，用于统计同时触发多次的状况，
     * 多次同时触发动画只显示一个，会增加num数量， 当每个请求结束时会减少num数量，当num为0时，结束动画
     * */

    if (config.showType === 'only') {
      div.className = 'isShow'
      div.setAttribute('data-num', '1')
    } else {
      div.setAttribute('data-num', '0')
    }
    document.body.appendChild(div)
    config.loadDOM = div
    callback ? callback() : ''
  }

  /**
   * 获取不同配置的节点
   * loading type，默认default，自定义image
   * @return html html节点
   * */
  _getTypeDom (type) {
    let html = ''
    if (type === 'default') {
      for (let i = 1; i <= 12; i++) {
        html += `<div class="loading-list loading-list-${i + 1}"></div>`
      }
    } else {
      html += `<div class='loading-main'></div>`
    }
    return html
  }

  /**
   * 添加css样式
   * @param config loading配置
   * @param callback 回调 可选
   * */
  _addStyle (config, callback) {
    this._loadCssCode(`#${config.loadingId}{
            position: ${config.width ? 'absolute' : 'fixed'};
            top: ${config.top ? `${config.top}px` : '0'};
            ${config.width ? `width:${config.width}px` : 'bottom:0'};
            ${config.height ? `height:${config.height}px` : ' right:0'};
            left: ${config.left ? `${config.left}px` : '0'};
            z-index: ${config.width ? '999' : '99999'} ;
            background:${config.background} ;
            background-size: 20%;
            display:none;
        }
        #${config.loadingId}.isShow{
            display:block;
        }
        #${config.loadingId} .m-loading-wrapper .loading-main{
                background:url(${config.loadingStyle.imgUrl}) no-repeat center;
                background-size:${config.loadingStyle.width};
                position: absolute;
                top: 0;
                bottom: 0;
                right: 0;
                left: 0;
            }`)
    callback ? callback() : ''
  }

  /**
   * 获取默认的style样式
   * */
  _getDefaultStyle () {
    let listCss = ''
    for (let i = 2; i <= 12; i++) {
      listCss += `
            .m-loading-wrapper .loading-list-${i} {
                -webkit-transform: rotate(${(i - 1) * 30}deg);
                transform: rotate(${(i - 1) * 30}deg);
            }
            .m-loading-wrapper .loading-list-${i}:before {
                -webkit-animation-delay: -${(12 + 1 - i) / 10}s;
                 animation-delay: -${(12 + 1 - i) / 10}s;
            }`
    }
    return `.m-loading-wrapper {
            width: 4em;
           height: 4em;
           position: absolute;
           top: 50%;
           left:50%;
           margin-left:-2em;
           margin-top:-2em;
        }
        .m-loading-wrapper .loading-list {
          width: 100%;
          height: 100%;
          position: absolute;
          left: 0;
          top: 0;
        }
        .m-loading-wrapper .loading-list:before {
          content: '';
          display: block;
          margin: 0 auto;
          width: 15%;
          height: 15%;
          background-color:#D10000;
          border-radius: 100%;
          -webkit-animation: m-loading-wrapper-delay 1.2s infinite ease-in-out both;
                  animation: m-loading-wrapper-delay 1.2s infinite ease-in-out both;
        }
        ${listCss}
        @-webkit-keyframes m-loading-wrapper-delay {
          0%, 80%, 100% {
            -webkit-transform: scale(0);
                    transform: scale(0);
          }
          40% {
            -webkit-transform: scale(1);
                    transform: scale(1);
          }
        }

        @keyframes m-loading-wrapper-delay {
          0%, 80%, 100% {
            -webkit-transform: scale(0);
                    transform: scale(0);
          }
          40% {
            -webkit-transform: scale(1);
                    transform: scale(1);
          }
        }`
  }

  /**
   * 动态添加css
   * @param code css样式
   * */
  _loadCssCode (code) {
    const style = document.createElement('style')
    style.type = 'text/css'
    style.rel = 'stylesheet'
    style.appendChild(document.createTextNode(code))
    const head = document.getElementsByTagName('head')[0]
    head.appendChild(style)
  }

  /**
   * 防止滚动穿透
   * */
  _fixedBody () {
    const scrollTop = document.body.scrollTop || document.documentElement.scrollTop
    document.body.style.cssText += `position:fixed;top:-${scrollTop}px`
  }

  /**
   * 滚动穿透恢复
   * */
  _looseBody () {
    const body = document.body
    body.style.position = ''
    const top = body.style.top
    document.body.scrollTop = document.documentElement.scrollTop = -parseInt(top)
    body.style.top = ''
  }

  /**
   * 获取距离页面顶部的距离
   * @param node 节点dom
   * */
  _getPositionTop (node) {
    let top = node.offsetTop
    let parent = node.offsetParent
    while (parent != null) {
      top += parent.offsetTop
      parent = parent.offsetParent
    }
    return top
  }

  /**
   * 获取距离页面左侧的距离
   * @param node 节点dom
   * */
  _getPositionLeft (node) {
    let left = node.offsetLeft
    let parent = node.offsetParent
    while (parent != null) {
      left += parent.offsetLeft
      parent = parent.offsetParent
    }
    return left
  }

  /**
   * 单独加载的loading
   * @param loadingConfigOnly 配置
   * */
  _onlyLoadingStart (loadingConfigOnly) {
    loadingConfigOnly = Object.assign({},
      this.config, {
        loadingId: `loading-${loadingConfigOnly.wrapperId}`,
        showType: 'only', // 默认全局加载，局部加载为only
      }, loadingConfigOnly)

    const loadingDomWrapper = document.getElementById(loadingConfigOnly.wrapperId)
    const loadingDom = document.getElementById(loadingConfigOnly.loadingId)

    const top = this._getPositionTop(loadingDomWrapper)
    const left = this._getPositionLeft(loadingDomWrapper)
    const width = loadingDomWrapper.offsetWidth
    const height = loadingDomWrapper.offsetHeight

    if (loadingDom) {
      loadingDom.className = 'isShow'
      const num = Number(loadingDom.getAttribute('data-num')) + 1
      loadingDom.setAttribute('data-num', num)
    } else {
      this._addDom(loadingConfigOnly)
      this._addStyle({
        ...loadingConfigOnly, top, left, width, height
      })
    }

    this.loadingSetTimeout[loadingConfigOnly.loadingId] = this.loadingSetTimeout[loadingConfigOnly.loadingId] || Date.parse(new Date())
  }

  /**
   * 单独加载的loading
   * @param loadingConfigOnly 配置
   * */
  _onlyLoadingEnd (loadingConfigOnly) {
    loadingConfigOnly = Object.assign({},
      this.config, {
        loadingId: `loading-${loadingConfigOnly.wrapperId}`,
        showType: 'only', // 默认全局加载，局部加载为only
      }, loadingConfigOnly)

    const loadingDom = document.getElementById(loadingConfigOnly.loadingId)

    const num = Number(loadingDom.getAttribute('data-num')) - 1
    loadingDom.setAttribute('data-num', num)
    if (num > 0) {
      return false
    }

    // 判断当前时间比初始的时间大多长
    const time = Date.parse(new Date()) - this.loadingSetTimeout[loadingConfigOnly.loadingId]
    console.log(loadingConfigOnly.minShowTime, this.config.minShowTime)
    const minShowTime = loadingConfigOnly.minShowTime || this.config.minShowTime
    if (time >= minShowTime) {
      loadingDom.className = ''
    } else {
      setTimeout(() => {
        loadingDom.className = ''
      }, minShowTime - time)
    }
  }

  /**
   * loading开始
   * @param loadingConfigOnly 可选，默认无为全屏显示
   * {
            wrapperId: 'top-wrapper',//loading显示地方的id (必填)
            loadingId: 'ceshi',
            loadingStyle: {
                type: 'image',
                width: '8em',
                imgUrl: 'https://assets.2dfire.com/frontend/56ef0094711c12fafd86e0a2ede914db.gif',
            },
            background: 'rgba(255,255,255,1)'
        }
   * @param loadingConfigOnly 单独loading的配置
   * */
  loadingStart (loadingConfigOnly) {
    // console.log(loadingConfigOnly)
    const that = this
    if (loadingConfigOnly.wrapperId) {
      this._onlyLoadingStart(loadingConfigOnly)
      return false
    }
    if (this.config.loadDOM) {
      this._fixedBody()
      this.config.loadDOM.className = 'isShow'
      this.num = Number(this.num) + 1
      this.config.loadDOM.setAttribute('data-num', this.num.toString())
      this.loadingSetTimeout[this.config.loadingId] = this.loadingSetTimeout[this.config.loadingId] || Date.parse(new Date())
    } else {
      this._addDom(this.config, () => {
        that.loadingStart(loadingConfigOnly)
      })
      this._addStyle(this.config)
    }
  }

  /**
   * loading结束
   * @param loadingConfigOnly 单独loading的配置
   * */
  loadingEnd (loadingConfigOnly) {
    // console.log(loadingConfigOnly)
    if (loadingConfigOnly.wrapperId) {
      this._onlyLoadingEnd(loadingConfigOnly)
      return false
    }
    if (!this.num) {
      return false
    }

    // 判断当前时间比初始的时间大多长
    const time = Date.parse(new Date()) - this.loadingSetTimeout[this.config.loadingId]
    // console.log(this.config.minShowTime)
    if (time >= this.config.minShowTime) {
      this._loadingDefaultDomSet()
    } else {
      setTimeout(() => {
        this._loadingDefaultDomSet()
      }, this.config.minShowTime - time)
    }
  }

  /**
   * loading结束时dom配置
   * */
  _loadingDefaultDomSet () {
    this.num = Number(this.num) - 1
    this.config.loadDOM.setAttribute('data-num', this.num.toString())
    if (this.num > 0) {
      return false
    }
    if (this.config.loadDOM && this.config.loadDOM.className.includes('isShow')) {
      this.config.loadDOM.className = ''
      this._looseBody()
    }
  }

  /**
   * loading强制结束
   * */
  mustEnd () {
    this.config.loadDOM.setAttribute('data-num', '0')
    this.config.loadDOM.className = ''
  }
}
