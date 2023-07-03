<template>
  <div class="page">
    <div class="gva-card-box">
      <div class="gva-card gva-top-card">
        <div class="gva-top-card-left">
          <div class="gva-top-card-left-title">任务信息</div>
          <!--<div class="gva-top-card-left-dot">今日晴，0℃ - 10℃，天气寒冷，注意添加衣物。</div>-->
          <div >
            <br><br>
            <el-row >
              <el-col :span="18" :xs="24" :sm="6" style="font-size: 88px">
                <div class="flex-center">
                  ID: {{task.id}}
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="12">
                <div class="flex-center">
                  任务名： {{task.name}}
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="6">
                <div class="flex-center">
                  状态：{{task.status}}
                </div>
              </el-col>
              <br><br><br><br>
              <el-col :span="18" :xs="24" :sm="6">
                <div class="flex-center">
                  创建者：{{task.creator}}
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="12">
                <div class="flex-center">
                  创建时间: {{dateFormatter(task.ct)}}
                </div>
              </el-col>
              <br><br><br><br>
              <el-col :span="18" :xs="24" :sm="12">
                <div class="flex-center">
                  任务描述: {{task.description}}
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
    </div>
    <div class="gva-card-box">
      <div class="gva-card">
        <div class="card-header">
          <div class="gva-top-card-left-title">任务项</div>
        </div>
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button size="small" icon="delete"  @click="openRejectDialog">驳回</el-button>
            <el-popover v-model:visible="execVisible" placement="top" width="160">
              <p>确定要执行吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button size="small" type="text" @click="execVisible = false">取消</el-button>
                <el-button size="small" type="primary" @click="onExec">确定</el-button>
              </div>
              <template #reference>
                <el-button icon="CaretRight" type="primary" size="small" style="margin-left: 10px;" @click="execVisible = true">执行</el-button>
              </template>
            </el-popover>
          </div>
          <el-table :data="tableData" style="width: calc(100% - 47px)" class="two-list">
            <el-table-column prop="id" label="序号"></el-table-column>
            <el-table-column prop="cmd" label="命令"></el-table-column>
            <el-table-column prop="cluster" label="集群">{{task.cluster}}</el-table-column>
            <el-table-column prop="db_name" label="库名">{{task.database}}</el-table-column>
            <el-table-column prop="exec_info" width="200" label="执行信息"></el-table-column>
            <el-table-column prop="cat_id" fixed="right" label="操作">
              <template #default="scope">
                <!-- <el-button
                    icon="edit"
                    size="small"
                    type="text"
                    @click="onExecAt(scope.row)"
                >此处开始</el-button> -->
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>

    <el-dialog v-model="rejectDialogFormVisible" :before-close="closeDialog" :title="驳回">
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="驳回原因" prop="description">
          <el-input
              v-model="form.reject_content"
              :autosize="{ minRows: 3, maxRows: 5000 }"
              type="textarea"
              placeholder="Please input"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="onReject">确 定</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { useRouter,useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import {
  getTask,
  updateTask,
} from '@/api/task/task'
import moment from "moment";
import {ElMessage} from "element-plus";

const form = ref({
  reject_content: ''
})

const taskType = 'redis'
const router = useRouter()
const rout = useRoute()
const execVisible = ref(false)

const dateFormatter = (seconds) =>{
  return moment( seconds *1000).format('YYYY-MM-DD HH:mm');
}

const newLineFormatter = (row, column) =>{
  return row.sql_content.replaceAll("\n", "<br>")
}

const task = ref({})
const tableData = ref([])

const getData = async() => {
  const response = await getTask(rout.params.id, taskType)
  task.value = response.data
  tableData.value = response.data.sub_tasks
}

getData()

const onReject = async() => {
  let params = {
    task:{
      id: Number(rout.params.id),
      reject_content: form.value.reject_content,
      action: 'reject',
      sub_task_type: 'redis',
    }
  }
  const res = await updateTask(params)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '操作成功',
      showClose: true
    })
  }
  rejectDialogFormVisible.value = false

  router.push({name: 'rReview'})
}

const onExec = async() => {
  execVisible.value = false
  let params = {
    task:{
      id: Number(rout.params.id),
      action: 'exec',
      sub_task_type: 'redis',
    }
  }
  const res = await updateTask(params)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '操作成功',
      showClose: true
    })
  }

  await syncUntilFinish()
}

const onExecAt = async(row) => {
  execVisible.value = false
  let params = {
    id: Number(rout.params.id),
    exec_item: row,
    action: 'beginAt'
  }
  const res = await updateTask(params)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '操作成功',
      showClose: true
    })
  }

  await syncUntilFinish()
}

const syncUntilFinish = async() =>{
  for (;;){
    await getData()
    if (task.value.status == 'execSuccess' || task.value.status == 'execFailed'){
      break
    }
  }
}

const rejectDialogFormVisible = ref(false)
const openRejectDialog = () => {
  rejectDialogFormVisible.value = true
}

const closeDialog = () =>{
  rejectDialogFormVisible.value = false
}

</script>
<script lang="ts">
export default {
  name: 'Dashboard'
}
</script>

<style lang="scss" scoped>
@mixin flex-center {
    display: flex;
    align-items: center;
}

.page {
    background: #f0f2f5;
    padding: 0;
    .gva-card-box{
      padding: 12px 16px;
      &+.gva-card-box{
        padding-top: 0px;
      }
    }
    .gva-card {
      box-sizing: border-box;
        background-color: #fff;
        border-radius: 2px;
        height: auto;
        padding: 26px 30px;
        overflow: hidden;
        box-shadow: 0 0 7px 1px rgba(0, 0, 0, 0.03);
    }
    .gva-top-card {
        height: 260px;
        @include flex-center;
        justify-content: space-between;
        color: #777;
        &-left {
          height: 100%;
          display: flex;
          flex-direction: column;
            &-title {
                font-size: 22px;
                color: #343844;
            }
            &-dot {
                font-size: 14px;
                color: #6B7687;
                margin-top: 24px;
            }
            &-rows {
                // margin-top: 15px;
                margin-top: 18px;
                color: #6B7687;
                width: 600px;
                align-items: center;
            }
            &-item{
              +.gva-top-card-left-item{
                margin-top: 24px;
              }
              margin-top: 14px;
            }
        }
        &-right {
            height: 600px;
            width: 600px;
            margin-top: 28px;
        }
    }
     ::v-deep(.el-card__header){
          padding:0;
          border-bottom: none;
        }
        .card-header{
          padding-bottom: 20px;
          border-bottom: 1px solid #e8e8e8;
        }
    .quick-entrance-title {
        height: 30px;
        font-size: 22px;
        color: #333;
        width: 100%;
        border-bottom: 1px solid #eee;
    }
    .quick-entrance-items {
        @include flex-center;
        justify-content: center;
        text-align: center;
        color: #333;
        .quick-entrance-item {
          padding: 16px 28px;
          margin-top: -16px;
          margin-bottom: -16px;
          border-radius: 4px;
          transition: all 0.2s;
          &:hover{
            box-shadow: 0px 0px 7px 0px rgba(217, 217, 217, 0.55);
          }
            cursor: pointer;
            height: auto;
            text-align: center;
            // align-items: center;
            &-icon {
                width: 50px;
                height: 50px !important;
                border-radius: 8px;
                @include flex-center;
                justify-content: center;
                margin: 0 auto;
                i {
                    font-size: 24px;
                }
            }
            p {
                margin-top: 10px;
            }
        }
    }
    .echart-box{
      padding: 14px;
    }
}
.dasboard-icon {
    font-size: 20px;
    color: rgb(85, 160, 248);
    width: 30px;
    height: 30px;
    margin-right: 10px;
    @include flex-center;
}
.flex-center {
    @include flex-center;
    font-size: 19px;
}

//小屏幕不显示右侧，将登陆框居中
@media (max-width: 750px) {
    .gva-card {
        padding: 20px 10px !important;
        .gva-top-card {
            height: auto;
            &-left {
                &-title {
                    font-size: 20px !important;
                }
                &-rows {
                    margin-top: 15px;
                    align-items: center;
                }
            }
            &-right {
                display: none;
            }
        }
        .gva-middle-card {
            &-item {
                line-height: 20px;
            }
        }
        .dasboard-icon {
            font-size: 18px;
        }
    }
}
</style>
