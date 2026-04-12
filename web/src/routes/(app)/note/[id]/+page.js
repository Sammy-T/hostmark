import { goto } from '$app/navigation';
import { error } from '@sveltejs/kit';

export const ssr = false

export async function load({ fetch, params }) {
    let resp = await fetch(`/api/note/${params.id}`);

    switch(resp.status) {
        case 200:
            break;

        case 401:
            const refResp = await fetch('/api/auth/refresh');

            switch(refResp.status) {
                case 200:
                    resp = await fetch(`/api/note/${params.id}`);
                    if(!resp.ok) error(resp.status, resp.statusText);
                    break;

                case 400:
                case 401:
                    goto('/login');

                default:
                    const msg = await refResp.text();
                    error(refResp.status, msg);
            }
            break;

        default:
            error(resp.status, resp.statusText);
    }

    const note = await resp.json();

    return { 
        note,
    };
}
