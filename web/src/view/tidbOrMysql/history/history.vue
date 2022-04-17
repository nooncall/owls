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
      <el-table :data="tableData" @sort-change="sortChange" row-key="id" @selection-change="handleSelectionChange">
        <el-table-column type="expand">
          <template #default="scope">
            <el-table :data="scope.row.exec_items" style="width: calc(100% - 47px)" class="two-list">
              <el-table-column prop="id" label="序号"></el-table-column>
              <el-table-column prop="cluster_name" label="集群"></el-table-column>
              <el-table-column prop="db_name" label="库名"></el-table-column>
              <el-table-column prop="task_type" label="类型"></el-table-column>
              <el-table-column prop="affect_rows" label="影响行数"></el-table-column>
              <el-table-column prop="status_name" label="状态"></el-table-column>
              <el-table-column prop="exec_info" label="执行信息"></el-table-column>
              <el-table-column class="cell" prop="cat_id" width="600" label="SQL语句">
                <template class="cell" style="white-space: pre-line;" #default="scope">
                  <code>
                    <span v-html="newLineFormatter(scope.row, '')"></span>
                  </code>
                </template>
              </el-table-column>
              <el-table-column prop="remark" label="备注" width="200"></el-table-column>
              <el-table-column prop="cat_id" fixed="right" label="操作">
                <template #default="scope">
                  <el-button
                      icon="edit"
                      size="small"
                      type="text"
                      v-if="scope.row.backup_status == 'backup_success'"
                      @click="rollbackFunc(scope.row)"
                  >回滚
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="150" prop="id" sortable="custom"/>
        <el-table-column align="left" label="任务名" min-width="150" prop="name" sortable="custom"/>
        <el-table-column align="left" label="状态" min-width="150" prop="status" sortable="custom"/>
        <el-table-column align="left" label="创建者" min-width="150" prop="creator" sortable="custom"/>
        <el-table-column align="left" label="创建时间" min-width="150" prop="ct" :formatter="dateFormatter"
                         sortable="custom"/>
        <el-table-column align="left" label="说明" min-width="150" prop="reject_content" sortable="custom"/>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 30, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="dialogTitle">
      <warning-bar title="原始数据如下,红色为修改列："/>

      <el-row>
        <el-col :span="12" v-for="column in rollbackColumn" style="font-weight: bold; font-size: 15px">
          {{ column }}
        </el-col>
      </el-row>
      <el-row v-for="item in rollbackData">
        <el-col :span="12" v-for="(field,idx) in item" v-bind:style="updatedIdx(idx)? 'color:red':''">
          {{ field }}
        </el-col>
      </el-row>

      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script lang="ts">
export default {
  name: 'Api',
}
</script>

<script lang="ts" setup>
import {
  listHistoryTask,} from '@/api/db/task'
import {toSQLLine} from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import {onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import moment from 'moment'
import {listBackup, rollback} from "../../../api/db/task";

const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const newLineFormatter = (row, column) => {
  return row.sql_content.replaceAll("\n", "<br>")
}

const apis = ref([])
const form = ref({
  cluster_name: '',
  db_name: '',
  task_type: '',
  remark: '',
  sql_content: ''
})
const methodOptions = ref([
  {
    value: 'POST',
    label: '创建',
    type: 'success'
  },
  {
    value: 'GET',
    label: '查看',
    type: ''
  },
  {
    value: 'PUT',
    label: '更新',
    type: 'warning'
  },
  {
    value: 'DELETE',
    label: '删除',
    type: 'danger'
  }
])

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

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const dateFormatter = (row, column) => {
  return moment(row.ct * 1000).format('YYYY-MM-DD HH:mm');
}

// 排序
const sortChange = ({prop, order}) => {
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
const getTableData = async () => {
  const table = await listHistoryTask({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
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

const deleteVisible = ref(false)

// 弹窗相关
const apiForm = ref(null)
const initForm = () => {
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const dialogTitle = ref('新增')
const dialogFormVisible = ref(false)
const openDialog = async (key) => {
  dialogTitle.value = '回滚'
  type.value = key
  dialogFormVisible.value = true
}
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}

const rollbackColumn = ref([])
const rollbackData = ref([])
const rollbackUpdateIdx = ref([])
const rollbackParam = ref({})

const rollbackFunc = async (row) => {
  // 获取回滚数据，展示及设置数据标志，确认执行回宫
  let params = {
    db_name: row.db_name,
    cluster_name: row.cluster_name,
    origin_sql: row.sql_content,
    backup_id: row.backup_id
  }
  rollbackParam.value = params

  let resp = await listBackup(params)
  rollbackData.value = resp.data.data_items
  rollbackColumn.value = resp.data.columns
  rollbackUpdateIdx.value = resp.data.index

  openDialog(params)
}

const updatedIdx = (idx) => {
  for (let i in rollbackUpdateIdx.value) {
    if (i == idx-1) {
      return true
    }
  }

  return false
}

const enterEditDialog = async () => {
  apiForm.value.validate(async valid => {
    if (valid) {
      switch (type.value) {
        case 'edit': {
          // 这里不传task id，后端判断是否包含id。
          // todo, refactor
          let params = {
            exec_item: form.value,
            action: "editItem"
          }
          const res = await updateTask(params)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '编辑成功',
              showClose: true
            })
          }
          getTableData()
          closeDialog()
        }
          break
        default:
          // eslint-disable-next-line no-lone-blocks
        {
          ElMessage({
            type: 'error',
            message: '未知操作',
            showClose: true
          })
        }
          break
      }
    }
  })
}


const enterDialog = async () => {
  const res = await rollback(rollbackParam.value)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '回滚成功',
      showClose: true
    })
  }
  getTableData()
  closeDialog()
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

.red {
  color: red;
}

.el-table .cell {
  white-space: pre-line;
}

</style>
