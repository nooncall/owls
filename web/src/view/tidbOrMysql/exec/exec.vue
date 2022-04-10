<template>
  <div class="page">
    <div class="gva-card-box">
      <div class="gva-card gva-top-card">
        <div class="gva-top-card-left">
          <div class="gva-top-card-left-title">任务信息</div>
          <!--<div class="gva-top-card-left-dot">今日晴，0℃ - 10℃，天气寒冷，注意添加衣物。</div>-->
          <div class="gva-top-card-left-rows">
          </div>
          <div >
            <el-row >
              <el-col :span="18" :xs="24" :sm="8" style="font-size: 88px">
                <div class="flex-center">
                  ID: 777
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="8">
                <div class="flex-center">
                  任务名： fix user type
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="8">
                <div class="flex-center">
                  状态：审核通过
                </div>
              </el-col>
              <br><br><br><br>
              <el-col :span="18" :xs="24" :sm="8">
                <div class="flex-center">
                  创建者：叶凡
                </div>
              </el-col>
              <el-col :span="18" :xs="24" :sm="8">
                <div class="flex-center">
                  创建时间：去年今日
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
            <el-button size="small" type="primary" icon="plus" @click="openDialog('add')">执行</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button size="small" type="text" @click="deleteVisible = false">取消</el-button>
                <el-button size="small" type="primary" @click="onDelete">确定</el-button>
              </div>
              <template #reference>
                <el-button icon="delete" size="small" style="margin-left: 10px;" @click="deleteVisible = true">驳回</el-button>
              </template>
            </el-popover>
          </div>
          <el-table :data="exec_items" style="width: calc(100% - 47px)" class="two-list">
            <el-table-column prop="id" label="序号"></el-table-column>
            <el-table-column prop="cluster_name" label="集群"></el-table-column>
            <el-table-column prop="db_name" label="库名"></el-table-column>
            <el-table-column prop="task_type" label="类型"></el-table-column>
            <el-table-column prop="affect_rows" label="影响行数"></el-table-column>
            <el-table-column prop="status" label="状态"></el-table-column>
            <el-table-column prop="exec_info" label="执行信息"></el-table-column>
            <el-table-column class="cell" prop="cat_id" width="600"  label="SQL语句">
              <template class="cell" style="white-space: pre-line;" #default="scope">
                <code>
                  <span v-html="newLineFormatter(scope.row, '')"></span>
                </code>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" width="200"></el-table-column>
            <el-table-column prop="cat_id" label="操作">
              <template #default="scope">
                <el-button
                    icon="edit"
                    size="small"
                    type="text"
                    @click="editTaskFunc(scope.row)"
                >从这里开始</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const toolCards = ref([
  {
    label: '用户管理',
    icon: 'monitor',
    name: 'user',
    color: '#ff9c6e',
    bg: 'rgba(255, 156, 110,.3)'
  },
  {
    label: '角色管理',
    icon: 'setting',
    name: 'authority',
    color: '#69c0ff',
    bg: 'rgba(105, 192, 255,.3)'
  },
  {
    label: '菜单管理',
    icon: 'menu',
    name: 'menu',
    color: '#b37feb',
    bg: 'rgba(179, 127, 235,.3)'
  },
  {
    label: '代码生成器',
    icon: 'cpu',
    name: 'autoCode',
    color: '#ffd666',
    bg: 'rgba(255, 214, 102,.3)'
  },
  {
    label: '表单生成器',
    icon: 'document-checked',
    name: 'formCreate',
    color: '#ff85c0',
    bg: 'rgba(255, 133, 192,.3)'
  },
  {
    label: '关于我们',
    icon: 'user',
    name: 'about',
    color: '#5cdbd3',
    bg: 'rgba(92, 219, 211,.3)'
  }
])

const router = useRouter()

const toTarget = (name) => {
  router.push({ name })
}

</script>
<script>
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
