<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="集群" prop="path">
          <el-autocomplete
              v-model="form.cluster_name"
              :fetch-suggestions="querySearchAsync"
              placeholder="请选择"
          />
        </el-form-item>
        <el-form-item label="库名" prop="apiGroup">
          <el-input v-model="form.db_name" placeholder="请输入，默认0"/>
        </el-form-item>
      </el-form>
    </div>

    <div>
      <div style="float: left; width: 55%">
        <el-input
            v-model="form.cmd"
            :autosize="{ minRows: 15, maxRows: 500 }"
            type="textarea"
            placeholder="Please input"
        />
        <div class="gva-btn-list" style="margin-top: 20px;float: right;">
          <el-button size="small" type="primary" icon="plus" @click="doReadData()">提交查询</el-button>
        </div>
      </div>

    </div>

    <div style="float: left;margin-left: 2%; width: 42%">
      <div class="gva-top-card-left-title">查询结果：</div>
      <br/>
      <code>
        <h4>{{ data }}</h4>
      </code>
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
  listClusterName,
} from '@/api/db/cluster'
import {
  readRedisData,
} from '@/api/db/read'

const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const apis = ref([])
const form = ref({
  cluster_name: '',
  db_name: 0,
  task_type: '',
  remark: '',
  cmd: ''
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
const data = ref([])
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
    cluster: form.value.cluster_name,
    db: parseInt(form.value.db_name),
    cmd: form.value.cmd,
  }
  const result = await readRedisData(params)
  if (result.code === 0) {
    data.value = result.data
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
  const resp = await listClusterName(true, 'redis')
  let result = []
  for (let cluster of resp.data) {
    result.push({value: cluster, name: cluster})
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
