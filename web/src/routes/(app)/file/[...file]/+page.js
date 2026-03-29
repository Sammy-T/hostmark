import { goto } from '$app/navigation';
import { error } from '@sveltejs/kit';
import { marked } from 'marked';

export const ssr = false

export async function load({ fetch, params }) {
    let resp = await fetch(`/api/file/${params.file}`);

    switch(resp.status) {
        case 200:
            break;

        case 400:
            return {};

        case 401:
            const refResp = await fetch('/api/auth/refresh');

            switch(refResp.status) {
                case 200:
                    resp = await fetch(`/api/file/${params.file}`);
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

    const markdown = await resp.text();
    const html = await marked.parse(markdown);

    return { 
        content: {
            markdown,
            html
        }
    };
}
