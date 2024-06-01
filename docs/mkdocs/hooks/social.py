""" Adds social params for each page. """

from mkdocs.plugins import event_priority
from mkdocs.structure.pages import Page
from mkdocs.config.defaults import MkDocsConfig


@event_priority(50)
def on_post_page(html: str, page: Page, config: MkDocsConfig) -> str:
    """
    This hook adds social meta tags for each page.
    """
    try:
        title = page.meta['social_meta']['title']
    except KeyError:
        title = page.title
    try:
        description = page.meta['social_meta']['description']
    except KeyError:
        description = config.site_description
    social_tags = config.theme.get_env().get_template('partials/meta_social_tags.html').render({
        'title': f'{config.theme["social_title_prefix"]} | {title}',
        'description': description,
        'url': page.canonical_url,
        'image': f'{config["site_url"]}{config.theme["social_image"]}',
    })
    return html.replace('</head>', f'{social_tags}</head>', 1)
