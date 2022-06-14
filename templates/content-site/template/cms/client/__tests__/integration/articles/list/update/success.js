import setup, { app } from '../../../helpers';
import axios from 'axios';
jest.mock( 'axios' );

beforeEach( () => {
    jest.resetModules();
    jest.clearAllMocks();
} );

it( 'success update article', async() => {
    expect.assertions( 7 );

    const DI = setup();
    axios
        .useDefaultMocks()
        .use( 'get', '/categories', {
            'categories': [
                { 'id': '1', 'name': 'Кухня', 'slug': 'kukhnia' },
                { 'id': '2', 'name': 'Быт', 'slug': 'byt' }
            ]
        } )
        .use( 'get', '/tags', {
            tags: [
                { 'id': '1', 'name': 'Hello world', 'slug': 'hello-world' },
                { 'id': '6', 'name': 'Название тега', 'slug': 'nazvanie-tega' }
            ]
        } )
        .use( 'get', '/articles', {
            'articles': [
                {
                    'id': '10',
                    'title': 'Тест2',
                    'slug': 'test2',
                    'category_id': '1',
                    'published_at': '2022-04-25T19:34:23.619+07:00'
                }
            ],
            'meta': { 'limit': 50, 'offset': 0, 'total': 1 }
        } )
        .use( 'get', '/articles/10', {
            'title': 'Тест2',
            'slug': 'test2',
            'category_id': '1',
            'short_body': 'Короткое',
            'body': '<p>Текст</p><p><u>Текст</u></p>',
            'published_at': '2022-04-25T19:34:23.619+07:00',
            'tags': [ 'hello-world' ]
        } )
        .use( 'put', '/articles/10', {} );

    history.pushState( {}, '', 'http://client.example.com/articles' );
    await app.start( DI );

    app.click( '[data-test-update-button="10"]' );
    await app.waitRendering();

    expect( window.location.href ).toBe( 'http://client.example.com/articles/10/edit' );

    app.checkMainPanel();

    const formData = {
        title: 'Travel',
        category_id: '2',
        short_body: 'Short body text',
        body: '<p>Full body text</p>',
        tags: [
            { id: '1', name: 'Hello world', slug: 'hello-world' },
            { id: '6', name: 'Название тега', slug: 'nazvanie-tega' },
            { name: 'New' }
        ],
        published_at: '2000-06-15T01:35:10.000Z'
    };

    app.form.fillBySelector( '[data-test-field-title] input', formData.title );
    app.form.fillBySelector( '[data-test-field-category] select', formData.category_id );
    app.form.fillBySelector( '[data-test-field-short-body] textarea', formData.short_body );
    app.form.fillBySelector( '[data-test-field-body] textarea', formData.body );

    const tagsSelector = '[data-test-field-tags] .autocomplete-input';
    formData.tags.forEach( tag => {
        const item = undefined === tag.slug ?
            'Add as new tag' :
            tag.name
        ;
        if ( null === document.querySelector( `[data-test-remove-tag="${item}"]` ) ) {
            app.form.fillBySelector( `${tagsSelector} input`, tag.name );
            app.click( `${tagsSelector} input` );
            app.click(
                `${tagsSelector} .autocomplete-results .result-item[data-test-item="${item}"]`
            );
        }
    } );

    const publishedAtSelector = '[data-test-field-published-at] .sham-ui-datetimepicker';

    app.click( `${publishedAtSelector} th.picker-switch` ); // select month
    app.click( `${publishedAtSelector} th.picker-switch` ); // select year
    app.click( `${publishedAtSelector} th.picker-switch` ); // select year decade
    app.click( `${publishedAtSelector} tbody tr td .decade:nth-child(1)` ); // select first decade
    app.click( `${publishedAtSelector} tbody tr td .year:nth-child(1)` ); // select 2000
    app.click( `${publishedAtSelector} tbody tr td .month:nth-child(6)` ); // select June
    app.click( `${publishedAtSelector} tbody tr:nth-child(3) .day:nth-child(4)` ); // select 15
    app.click( `${publishedAtSelector} tfoot tr th` ); // select time
    app.click( `${publishedAtSelector} tbody tr .hour` ); // select hour
    app.click( `${publishedAtSelector} tbody tr td span:nth-child(10)` ); // select 09
    app.click( `${publishedAtSelector} tbody tr .minute` ); // select minutes
    app.click( `${publishedAtSelector} tbody tr td span:nth-child(8)` ); // select 35
    app.click( `${publishedAtSelector} tbody tr .second` ); // select seconds
    app.click( `${publishedAtSelector} tbody tr td span:nth-child(3)` ); // select 10

    app.checkMainPanel();

    app.form.submit();

    app.click( '[data-test-modal] [data-test-ok-button]' );
    await app.waitRendering();
    app.checkMainPanel();

    expect( axios.mocks.put ).toHaveBeenCalledTimes( 1 );
    expect( axios.mocks.put.mock.calls[ 0 ][ 0 ] ).toBe( '/articles/10' );
    expect( axios.mocks.put.mock.calls[ 0 ][ 1 ] ).toEqual( formData );
} );

