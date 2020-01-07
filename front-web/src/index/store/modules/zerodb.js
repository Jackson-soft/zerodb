import API from "@api"

const zerodb = {
    namespaced: true,
    state: {
        snapshotList: [],
        clusterList: []
    },
    mutations: {
        _getSnapshotList(state, data) {
            state.snapshotList = data
        },
        _getClusterList(state, data) {
            state.clusterList = data
        }
    },
    actions: {
        getSnapshotList(context, clusterName) {
            API.SnapshotList({ clusterName }).then(res => {
                context.commit("_getSnapshotList", res)
            })
        },
        getClusterList(context) {
            API.ClusterList().then(res => {
                context.commit("_getClusterList", res)
            })
        }
    },
    getters: {
        snapshotList(state) {
            return state.snapshotList
        },
        clusterList(state) {
            return state.clusterList
        }
    }
}
export default zerodb
