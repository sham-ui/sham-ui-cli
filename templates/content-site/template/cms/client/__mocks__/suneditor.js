class MockEditor {
    static create( el ) {
        return new MockEditor( el );
    }

    constructor( el ) {

        // Use hack for correct display in snapshots
        const text = document.createTextNode( el.value );
        el.appendChild( text );
        this.el = el;
        this.textNode = text;
        this._onChangeHandler = () => {
            this.textNode.textContent = el.value;
            this.onChange( el.value );
        };
        el.addEventListener( 'change', this._onChangeHandler );
    }

    setContents( content ) {
        this.el.value = content;
        this.textNode.textContent = content || '';
    }

    getContents() {
        return this.el.value;
    }

    destroy() {
        this.el.removeEventListener( 'change', this._onChangeHandler );
        this.el = null;
        this.textNode = null;
    }
}

export default {
    init() {
        return MockEditor;
    }
};
