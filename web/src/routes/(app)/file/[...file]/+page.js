import { error } from '@sveltejs/kit';
import { marked } from 'marked';

export const ssr = false

export async function load({ fetch, params }) {
    const resp = await fetch(`/api/file/${params.file}`);
    if(!resp.ok) {
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
