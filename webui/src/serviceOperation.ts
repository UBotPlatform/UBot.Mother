import {motherBaseURL} from './buildConfig'
export enum ServiceStatus {
    Stopped,
    Running,
    Starting,
    Exited
}
export interface ServiceInfo {
    id: string
    status: ServiceStatus
    launch_at: string
}
export async function startService(info: ServiceInfo) {
    fetch(motherBaseURL + "/api/mother/services", {
        body: JSON.stringify([
            {
                "id": info.id, 
                "status": ServiceStatus.Starting
            }
        ]),
        headers: {
            'content-type': 'application/json'
        },
        method: 'PUT'
    })
}
export async function stopService(info: ServiceInfo) {
    fetch(motherBaseURL + "/api/mother/services", {
        body: JSON.stringify([
            {
                "id": info.id, 
                "status": ServiceStatus.Stopped
            }
        ]),
        headers: {
            'content-type': 'application/json'
        },
        method: 'PUT'
    })
}