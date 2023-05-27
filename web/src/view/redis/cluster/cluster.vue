<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模糊查询">
          <el-input v-model="searchInfo.key" placeholder="模糊查询" />
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button size="small" icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="small" type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="small" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="small" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="delete" size="small" :disabled="!apis.length" style="margin-left: 10px;" @click="deleteVisible = true">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
            type="selection"
            width="55"
        />
        <el-table-column align="left" label="集群名" min-width="150" prop="name" sortable="custom" />
        <el-table-column align="left" label="描述" min-width="150" prop="description" sortable="custom" />
        <el-table-column align="left" label="地址" min-width="150" prop="addr" sortable="custom" />
        <el-table-column align="left" label="用户名" min-width="150" prop="user" sortable="custom" />
        <el-table-column align="left" label="创建时间" min-width="150" prop="ct" sortable="custom" />
        <el-table-column align="left" label="更新时间" min-width="150" prop="ut" sortable="custom">
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFiletr(scope.row.method) }}
            </div>
          </template>
        </el-table-column>

        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
                icon="edit"
                size="small"
                type="text"
                @click="editClusterFunc(scope.row)"
            >编辑</el-button>
            <el-button
                icon="delete"
                size="small"
                type="text"
                @click="deleteClusterFunc(scope.row)"
            >删除</el-button>
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

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="dialogTitle">
      <warning-bar title="新增集群，需要在角色管理内配置权限才可使用" />
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="path">
          <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="地址" prop="path">
          <el-input v-model="form.addr" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" prop="apiGroup">
          <el-input v-model="form.description" autocomplete="off" />
        </el-form-item>
        <el-form-item label="用户名" prop="description">
          <el-input v-model="form.user" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="description">
          <el-input v-model="form.pwd" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Api',
}
</script>

<script setup>
import {
  listCluster,
  createCluster,
  updateCluster,
  deleteCluster,
} from '@/api/db/cluster'
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
  c_type: 'redis',
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
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  addr: [
    { required: true, message: '请输入地址', trigger: 'blur' }
  ],
  pwd: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  user: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
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
  const table = await listCluster({ page: page.value, 
    pageSize: pageSize.value, 
    type: "redis",
    ...searchInfo.value })
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
const onDelete = async() => {
  const ids = apis.value.map(item => item.ID)
  const res = await deleteApisByIds({ ids })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
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

const dialogTitle = ref('新增')
const dialogFormVisible = ref(false)
const openDialog = (key) => {
  switch (key) {
    case 'add':
      dialogTitle.value = '新增'
      break
    case 'edit':
      dialogTitle.value = '编辑'
      break
    default:
      break
  }
  type.value = key
  dialogFormVisible.value = true
}
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}

const editClusterFunc = async(row) => {
  form.value = row
  openDialog('edit')
}

const enterDialog = async() => {
  apiForm.value.validate(async valid => {
    if (valid) {
      switch (type.value) {
        case 'add':
        {
          const res = await createCluster(form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功',
              showClose: true
            })
          }
          getTableData()
          closeDialog()
        }

          break
        case 'edit':
        {
          const res = await updateCluster(form.value)
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

const deleteClusterFunc = async(row) => {
  ElMessageBox.confirm('确定删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
      .then(async() => {
        const res = await deleteCluster(row)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '删除成功!'
          })
          if (tableData.value.length === 1 && page.value > 1) {
            page.value--
          }
          getTableData()
        }
      })
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
