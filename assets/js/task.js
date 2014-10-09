Vue.component('task-list', {
    methods: {
        getTask: function (id) {
            var task = null
            this.data.forEach(function (d) {
                d.tasks.forEach(function (t) {
                    if (t.id === id) task = t
                })
            })
            return task
        },
        getProject: function (id, name) {
            var project = null
            this.data.forEach(function (d) {
                if (d.id === id) project = d
            })
            if (project === null) {
                project = {
                    id: id,
                    name: name,
                    tasks: []
                }
                this.data.unshift(project)
            }
            return project
        },
        refreshTask: function (pid, tid, task) {
            if (pid == null) {
                pid = this.getTask(tid).project_id
            }
            var project = this.getProject(pid)
            var index = -1
            for (var i=0; i<project.tasks.length; i++) {
                if (project.tasks[i].id === tid) {
                    index = i
                    break
                }
            }
            if (index >= 0) project.tasks.splice(index, 1)
            if (task) project.tasks.unshift(task)
        },
        report: function (id) {
            var task = this.getTask(id)
            var model = {
                project_id: task.project_id,
                progress: task.progress,
                update_at: task.update_at,
                task_id: task.id
            }

            vui.openbox({
                title: task.name,
                src: 'report.list',
                data: {
                    model: model
                },
                show: true,
                callback: function () {
                    if (task.update_at != model.update_at) {
                        task.progress = model.progress
                        this.refreshTask(task.project_id, task.id, task)
                    }
                }.bind(this)
            })
        },
        remove: function (id, op) {
            var self = this
            vui.openbox.confirm("确定" + op + "这个任务？", function (status) {
                if (!status) return

                self.act('task/', id, 'del', '任务' + op + '成功')
            })
        },
        act: function (src, id, method, msg) {
            var self = this
            vui.loading.start()
            vui.request[method](src).send(id).end(function (res) {
                vui.loading.end()

                if (res.status !== 200) {
                    vui.message.error('', res.status)
                    return
                }
                if (res.body.status === 0) {
                    vui.message.error(res.body.errors || res.body.msg)
                    return
                }

                self.refreshTask(res.body.data, id)
                if (msg) vui.message.info(msg)
            }) 
        },
        finish: function (id) {
            this.act('task/finish', id, 'post')
        },
        taskEdit: function (id, name) {
            var model = this.getTask(id)
            var box = vui.openbox({
                title: name,
                show: true,
                width: 8,
                src: 'task_edit',
                data: { 
                    model: vui.utils.copy(model)
                },
                callback: function (model) {
                    if (!model) return
                    this.refreshTask(model.project_id, model.id, model)
                }.bind(this)
            })
        }
    },
    data: {
        data: []
    },
    created: function () {
        vui.loading.start()
        vui.request.get('task/').end(function (res) {
            vui.loading.end()
            if (res.status !== 200 || res.body.status !== 1) {
                vui.message.error(res.body.errors, res.status)
                return
            }

            this.data = res.body.data
        }.bind(this))
    }
})
