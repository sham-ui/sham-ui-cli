class MockEditor {
    static create( el ) {
        return new MockEditor( el );
    }

    constructor( el ) {
        this.el = el;
        this._onChangeHandler = () => this.onChange( el.value );
        el.addEventListener( 'change', this._onChangeHandler );
    }

    destroy() {
        this.el.removeEventListener( 'change', this._onChangeHandler );
        this.el = null;
    }
}

export default {
    init() {
        return MockEditor;
    }
};
