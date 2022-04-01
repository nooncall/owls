<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模糊查询">
          <el-input v-model="searchInfo.key" placeholder="模糊查询"/>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button size="small" icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
            type="selection"
            width="55"
        />
        <el-table-column align="left" label="名称" min-width="150" prop="name" sortable="custom"/>
        <el-table-column align="left" label="规则" min-width="150" prop="summary" sortable="custom"/>

        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-switch
                v-model="scope.row.open"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @click="onSwitch(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Api',
}
</script>

<script setup>
import {
  listRule,
  updateRuleStatus,
} from '@/api/db/rule'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const apis = ref([])
const form = ref({
  name: '',
  addr: '',
  pwd: '',
  user: '',
})
const methodOptions = ref([
  {
    value: 'GET',
    label: '查看',
    type: '',
  },
  {
    value: 'PUT',
    label: '更新',
    type: 'warning',
  },
])

const type = ref('')

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
}
// 搜索

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

const onSwitch = async(row) => {
  let action = 'open'
  if (!row.open) {
    action = 'close'
  }
  await updateRuleStatus({ name: row.name, action: action })
  await getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 排序
const sortChange = ({ prop, order }) => {
  if (prop) {
    if (prop === 'ID') {
      prop = 'id'
    }
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await listRule({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// 批量操作
const handleSelectionChange = (val) => {
  apis.value = val
}

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
</style>
