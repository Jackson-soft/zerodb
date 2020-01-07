<template>
  <div>
    <Table ref="table" :columns="columns" v-bind="$props">
    </Table>
    <div style="display:none">
      <slot></slot>
    </div>
  </div>
</template>
<script>
export default {
  provide () {
    // const self = this
    return {
      updateTable () {
        const parent = this.$parent // 防抖
        if (parent.throttle) {
          return
        }
        parent.throttle = true
        setTimeout(() => {
          parent.getColumn()
          parent.throttle = false
        }, 100)
        // this.$nextTick(() => {
        //   self.getColumn()
        // })
      }
    }
  },
  props: {
    data: Array,
    stripe: Boolean,
    border: Boolean,
    showHeader: {
      type: Boolean,
      default: true
    },
    width: [Number, String],
    height: [Number, String],
    loading: Boolean,
    disabledHover: Boolean,
    highlightRow: Boolean,
    rowClassName: Function,
    size: String,
    noDataText: String,
    noFilteredDataText: String,
    refresh: [Number, String]
  },
  data () {
    return {
      columns: [],
      throttle: false
    }
  },
  methods: {
    getColumn () {
      const columnsArr = []
      if (this.$slots.default) {
        const slots = Array.isArray(this.$slots.default)
          ? this.$slots.default
          : [this.$slots.default]
        Array.from(slots)
          .filter(component => component.componentInstance)
          .forEach(column => {
            const instance = column.componentInstance
            const o = instance.$options.propsData
            o.key = o.prop
            if (instance.$scopedSlots.default) {
              o.render = (h, props) => h(
                  'div',
                  instance.$scopedSlots.default(props.row)
                )
            }
            columnsArr.push(o)
          })
        this.columns = columnsArr
      }
    },
    methodsInit () {
      [
        'on-current-change',
        'on-select',
        'on-select-cancel',
        'on-select-all',
        'on-selection-change',
        'on-sort-change',
        'on-filter-change',
        'on-row-click',
        'on-row-dblclick',
        'on-expand'
      ].forEach(evtName => {
        this.$refs.table.$on(evtName, (...args) => {
          this.$emit(evtName, ...args)
        })
      });

      ['exportCsv', 'clearCurrentRow'].forEach(methodName => {
        this[methodName] = this.$refs.table[methodName]
      })
    }
  },
  mounted () {
    this.getColumn()
    this.methodsInit()
  }
}
</script>
