import React from 'react';

function pageHref(num) {
    return `#page-${num + 1}`
}

export class Paginator extends React.Component {
    constructor(props, context) {
        super(props, context);

        this.state = {
            paginationFixed: true
        };
    }
    refixPosition() {
        if ($("#output").height() < $(window).height() && this.state.paginationFixed) {
            console.log("Relative position:", $("#output").height(), " > ", $(window).height());
            this.setState({paginationFixed: false});
        } else if ($("#output").height() >= $(window).height() && !this.state.paginationFixed) {
            console.log("paginationFixed: true")
            this.setState({paginationFixed: true});
        }
    }

    componentDidMount() {
        this.refixPosition();
    }

    componentDidUpdate() {
        this.refixPosition();
    }

    handlePrevious(e) {
        e.preventDefault()
        this.props.onPageChange(this.props.currentPage - 1)
    }

    handleNext(e) {
        e.preventDefault()
        this.props.onPageChange(this.props.currentPage + 1);
    }

    handlePageButton(page, e) {
        e.preventDefault();
        this.props.onPageChange(page);
    }

    handleSetCurrentPage(selectedPageNumber, e) {
        if (this.props.numPages > selectedPageNumber || selectedPageNumber < 1) {
            console.log("Not possible ERRRR")
        } else {
            this.props.onPageChange(number);
        }
    }

    renderPrevious() {
        if(this.props.currentPage <= 0) {
            return (<a className='reactable-previous-page disabled'>Previous</a>);
        }

        return (
            <a  className='reactable-previous-page'
                onClick={this.handlePrevious.bind(this)}
                href={pageHref(this.props.currentPage - 1)}>
                Previous
           </a>
        );
    }

    renderNext() {
        if(this.props.currentPage >= this.props.numPages - 1) {
            return (<a className='reactable-next-page disabled'>Next</a>);
        }

        return (
            <a  className='reactable-next-page'
                onClick={this.handleNext.bind(this)}
                href={pageHref(this.props.currentPage + 1)}>
                Next
           </a>
        );
    }

    renderPageButton() {
        return (<span className='reactable-page'> {this.props.currentPage + 1} out of {this.props.numPages}</span>);
    }

    render() {
        if (typeof this.props.colSpan === 'undefined') {
            throw new TypeError('Must pass a colSpan argument to Paginator');
        }

        if (typeof this.props.numPages === 'undefined') {
            throw new TypeError('Must pass a non-zero numPages argument to Paginator');
        }

        if (typeof this.props.currentPage === 'undefined') {
            throw new TypeError('Must pass a currentPage argument to Paginator');
        }

        // let pageButtons = [];
        // let pageButtonLimit = this.props.pageButtonLimit;
        // let currentPage = this.props.currentPage;
        // let numPages = this.props.numPages;
        // let lowerHalf = Math.round( pageButtonLimit / 2 );
        // let upperHalf = (pageButtonLimit - lowerHalf);

        // for (let i = 0; i < this.props.numPages; i++) {
        //     let showPageButton = false;
        //     let pageNum = i;
        //     let className = "reactable-page-button";
        //     if (currentPage === i) {
        //         className += " reactable-current-page";
        //     }
        //     pageButtons.push( this.renderPageButton(className, pageNum));
        // }

        // if(currentPage - pageButtonLimit + lowerHalf > 0) {
        //     if(currentPage > numPages - lowerHalf) {
        //         pageButtons.splice(0, numPages - pageButtonLimit)
        //     } else {
        //         pageButtons.splice(0, currentPage - pageButtonLimit + lowerHalf);
        //     }
        // }

        // if((numPages - currentPage) > upperHalf) {
        //     pageButtons.splice(pageButtonLimit, pageButtons.length - pageButtonLimit);
        // }

         // <tbody className="reactable-pagination">
         //        <tr>
         //            <td colSpan={this.props.colSpan}>
         //                <div className="reactable-pagination-button-container">
         //                {this.renderPrevious()}
         //                {this.renderPageButton()}
         //                {this.renderNext()}
         //                </div>
         //            </td>
         //        </tr>
         //    </tbody>

        if (this.state.paginationFixed) {
            var style = { position: 'fixed' }
            $("#results").attr('margin-bottom', '21px');
        } else {
            var style = { position: 'relative' }
            $("#results").attr('margin-bottom', '0');
        }

        return (
            <div className="reactable-pagination">
                <div className="reactable-pagination-button-container" style={style}>
                    {this.renderPrevious()}
                    {this.renderPageButton()}
                    {this.renderNext()}
                </div>
            </div>
        );
    }
};

