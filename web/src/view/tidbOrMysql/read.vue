<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="集群" prop="path">
          <el-autocomplete
              v-model="form.cluster_name"
              :fetch-suggestions="querySearchAsync"
              placeholder="请选择"
              @select="handleSelect"
          />
        </el-form-item>
        <el-form-item label="库名" prop="apiGroup">
          <el-autocomplete
              v-model="form.db_name"
              :fetch-suggestions="querySearchDBAsync"
              placeholder="请选择"
              @select="handleDBSelect"
          />
        </el-form-item>
        <el-form-item label="表名" prop="apiGroup">
          <el-autocomplete
              v-model="form.table_name"
              :fetch-suggestions="querySearchTableAsync"
              placeholder="请选择"
              @select="handleTableSelect"
          />
        </el-form-item>
      </el-form>
    </div>

    <div>
      <div style="float: left; width: 46%">
        <el-input
            v-model="form.sql_content"
            :autosize="{ minRows: 15, maxRows: 500 }"
            type="textarea"
            placeholder="Please input"
        />
        <div class="gva-btn-list" style="margin-top: 20px;float: right;">
          <el-button size="small" type="primary" icon="plus" @click="doReadData()">提交查询</el-button>
        </div>
      </div>

      <div style="float: left;margin-left: 3%; width: 45%">
        <code>
          <span v-html="tableInfoFormatted()"></span>
        </code>
      </div>

    </div>

    <br><br>

    <div class="gva-table-box" style="margin-top: 20px;float: left; width: 95%">
      <el-table :data="tableData" >
        <el-table-column min-width="160" :label="date" v-for="(date, key) in column" :render-header="renderHeader">
          <template align="left" #default="scope">
            {{tableData[scope.$index][key]}}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  name: 'Api',
}
</script>

<script lang="ts" setup>
import {toSQLLine} from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import {onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import moment from 'moment'
import {
  listCluster,
  listDatabase,
  listTable,
} from '@/api/db/cluster'
import {
  readData,
  getTableInfo,
} from '@/api/db/read'

const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const tableInfoFormatted = () => {
  return tableInfo.value.replaceAll("\n", "<br><br>")
}

const apis = ref([])
const form = ref({
  cluster_name: '',
  db_name: '',
  task_type: '',
  remark: '',
  sql_content: ''
})

const type = ref('')
const rules = ref({
  name: [{required: true, message: '请输入名称', trigger: 'blur'}],
  addr: [
    {required: true, message: '请输入地址', trigger: 'blur'}
  ],
  pwd: [
    {required: true, message: '请输入密码', trigger: 'blur'}
  ],
  user: [
    {required: true, message: '请输入用户名', trigger: 'blur'}
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const column = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
}

const dateFormatter = (row, column) => {
  return moment(row.ct * 1000).format('YYYY-MM-DD HH:mm');
}

const renderHeader = ({column,index}) =>{
  let l = column.label.length;
  let f = 12; //每个字的比例值，大概会比字体大小大一点
  column.minWidth = f*l < 150? 150: f*l; //字大小乘个数即长度,都是比例值
  return ('span',{class:'gva-table-box',style:{width:'100%'}},[column.label]);
}

// 查询
const doReadData = async () => {
  let params = {
    cluster_name: form.value.cluster_name,
    db_name: form.value.db_name,
    table_name: form.value.table_name,
    sql: form.value.sql_content,
  }
  const table = await readData(params)
  if (table.code === 0) {
    tableData.value = table.data.data_items
    column.value= table.data.columns
  }
}

// 批量操作
const handleSelectionChange = (val) => {
  apis.value = val
}

// 弹窗相关
const apiForm = ref(null)
const initForm = () => {
  apiForm.value.resetFields()
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const state = ref('')

interface clusterItem {
  name: string
  value: string
}

const clusters = ref<clusterItem[]>([])
const db = ref<clusterItem[]>([])
const table = ref<clusterItem[]>([])

const loadCluster = async () => {
  const resp = await listCluster({page: 1, pageSize: 5})
  let result = []
  for (let cluster of resp.data.list) {
    result.push({value: cluster.name, name: cluster.name})
  }

  return result
}

const loadDB = async (cluster) => {
  const resp = await listDatabase(cluster)
  let result = []
  for (let db of resp.data) {
    result.push({value: db, name: db})
  }

  return result
}

const loadTable = async (cluster, db) => {
  const resp = await listTable(cluster, db)
  let result = []
  for (let db of resp.data) {
    result.push({value: db, name: db})
  }

  return result
}

let timeout: NodeJS.Timeout
const querySearchAsync = (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
      ? clusters.value.filter(createFilter(queryString))
      : clusters.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 3000 * Math.random())
}

const querySearchDBAsync = async (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
      ? db.value.filter(createFilter(queryString))
      : db.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 3000 * Math.random())
}

const querySearchTableAsync = async (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
      ? table.value.filter(createFilter(queryString))
      : table.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 3000 * Math.random())
}

const createFilter = (queryString: string) => {
  return (restaurant: clusterItem) => {
    return (
        restaurant.name.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

const handleSelect = async (item: clusterItem) => {
  db.value = await loadDB(item.value)
}

const handleDBSelect = async (item: clusterItem) => {
  table.value = await loadTable(form.value.cluster_name, item.value)
}

const tableInfo = ref('')
const handleTableSelect = async (item: clusterItem) => {
  let params = {
    cluster_name: form.value.cluster_name,
    db_name: form.value.db_name,
    table_name: form.value.table_name,
  }
  let result = await getTableInfo(params)
  tableInfo.value = result.data
}

onMounted(async () => {
  clusters.value = await loadCluster()
})

</script>


<style scoped lang="scss">
.button-box {
  padding: 10px 20px;

  .el-button {
    float: right;
  }
}

.warning {
  color: #dc143c;
}

.el-table .cell {
  white-space: pre-line;
}

</style>
