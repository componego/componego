""" Adds additional logic to improve navigation. """

from mkdocs.plugins import event_priority
from mkdocs.structure.pages import Page
from mkdocs.config.defaults import MkDocsConfig
from bs4 import BeautifulSoup


@event_priority(50)
def on_post_page(html: str, page: Page, config: MkDocsConfig) -> str:
    """
    This hook changes the HTML tree by adding or removing some nodes or attributes.
    """
    parsed_html = BeautifulSoup(html, 'html.parser')
    items = parsed_html.select('ul.md-nav__list .md-nav__item--section')
    # Menu items "Built-in Components" and "Examples".
    for index in [2, 3]:
        # The necessary checks for the keys in the list are missing
        # because we expect an exception if such a key does not exist.
        section = items[index]
        section['class'].remove('md-nav__item--section')
        section['class'] += ['toggle-color']
        for toggle in section.select('.md-toggle--indeterminate'):
            toggle['class'].remove('md-toggle--indeterminate')
        for link in section.select('a'):
            if 'https://github.com' in link.get('href', ''):
                link['target'] = '_blank'
    for link in parsed_html.select('a[target="_blank"]'):
        if link.get('rel', None) is None:
            # noinspection SpellCheckingInspection
            link['rel'] = 'noopener'
    return str(parsed_html)
