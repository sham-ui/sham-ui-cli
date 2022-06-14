import * as directives from 'sham-ui-directives';
import { onFocusIn } from '../../../../../../../src/directives/on-focus-in';
// eslint-disable-next-line max-len
import RoutesArticlesCommonFieldsTags  from '../../../../../../../src/components/routes/articles/common/fields/tags.sfc';
import renderer from 'sham-ui-test-helpers';

it( 'renders correctly', () => {
    const meta = renderer( RoutesArticlesCommonFieldsTags, {
        all: [
            { name: 'first' },
            { name: 'second' }
        ],
        directives: {
            onFocusIn,
            ...directives
        }
    } );
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'show all items on click', () => {
    const meta = renderer( RoutesArticlesCommonFieldsTags, {
        all: [
            { name: 'first' },
            { name: 'second' }
        ],
        directives: {
            onFocusIn,
            ...directives
        }
    } );
    document.querySelector( 'input' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
} );

it( 'add tag', () => {
    const onChange = jest.fn();
    const tags = [
        { name: 'first' },
        { name: 'second' }
    ];
    const meta = renderer( RoutesArticlesCommonFieldsTags, {
        all: tags,
        selected: [],
        onChange,
        directives: {
            onFocusIn,
            ...directives
        }
    } );
    document.querySelector( 'input' ).click();
    document.querySelector( '[data-test-item="second"]' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
    expect( onChange ).toHaveBeenCalledTimes( 1 );
    expect( onChange.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( onChange.mock.calls[ 0 ][ 0 ] ).toEqual( [ tags[ 1 ] ] );
} );


it( 'remove tag', () => {
    const onChange = jest.fn();
    const tags = [
        { name: 'first' },
        { name: 'second' }
    ];
    const meta = renderer( RoutesArticlesCommonFieldsTags, {
        all: tags,
        selected: [ ...tags ],
        onChange,
        directives: {
            onFocusIn,
            ...directives
        }
    } );
    document.querySelector( '[data-test-remove-tag="second"]' ).click();
    expect( meta.toJSON() ).toMatchSnapshot();
    expect( onChange ).toHaveBeenCalledTimes( 1 );
    expect( onChange.mock.calls[ 0 ] ).toHaveLength( 1 );
    expect( onChange.mock.calls[ 0 ][ 0 ] ).toEqual( [ tags[ 0 ] ] );
} );
