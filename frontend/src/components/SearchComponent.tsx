import React, { useState, useEffect } from 'react';
import { Search, X, Clock, TrendingUp } from 'lucide-react';

interface SearchResult {
  id: string;
  title: string;
  url: string;
  description: string;
  type: 'web' | 'image' | 'news';
}

interface RecentSearch {
  id: string;
  query: string;
  timestamp: number;
}

interface TrendingTopic {
  id: string;
  title: string;
  searches: number;
}

const SearchComponent: React.FC = () => {
  const [query, setQuery] = useState<string>('');
  const [isSearching, setIsSearching] = useState<boolean>(false);
  const [showSuggestions, setShowSuggestions] = useState<boolean>(false);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [recentSearches] = useState<RecentSearch[]>([
    { id: '1', query: 'MacBook Pro M3', timestamp: Date.now() - 3600000 },
    { id: '2', query: 'iPhone 15 features', timestamp: Date.now() - 7200000 },
    { id: '3', query: 'Apple Vision Pro', timestamp: Date.now() - 10800000 },
  ]);
  const [trendingTopics] = useState<TrendingTopic[]>([
    { id: '1', title: 'WWDC 2024', searches: 2540000 },
    { id: '2', title: 'iOS 18 updates', searches: 1890000 },
    { id: '3', title: 'Apple Intelligence', searches: 1650000 },
    { id: '4', title: 'M4 chip performance', searches: 1420000 },
  ]);

  const handleSearch = async (searchQuery: string = query) => {
    if (!searchQuery.trim()) return;
    
    setIsSearching(true);
    setShowSuggestions(false);
    
    try {
      const response = await fetch(`/api/search?q=${encodeURIComponent(searchQuery)}`);
      const data = await response.json();
      setSearchResults(data.results);
    } catch (error) {
      console.error('Search failed:', error);
    } finally {
      setIsSearching(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    setShowSuggestions(value.length > 0);
    
    if (value.length === 0) {
      setSearchResults([]);
    }
  };

  const clearSearch = () => {
    setQuery('');
    setSearchResults([]);
    setShowSuggestions(false);
  };

  const handleSuggestionClick = (suggestion: string) => {
    setQuery(suggestion);
    handleSearch(suggestion);
  };

  useEffect(() => {
    const handleClickOutside = () => {
      setShowSuggestions(false);
    };

    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-white">
      {/* Header */}
      <header className="backdrop-blur-xl bg-white/80 border-b border-gray-200/50 sticky top-0 z-50">
        <div className="max-w-4xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Search className="w-4 h-4 text-white" />
              </div>
              <h1 className="text-xl font-semibold text-gray-900">Search</h1>
            </div>
            <div className="text-sm text-gray-500">
              Powered by Intelligence
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-12">
        {/* Search Section */}
        <div className="relative mb-12">
          <div className="relative">
            <div className="relative group">
              <input
                type="text"
                value={query}
                onChange={handleInputChange}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                onClick={(e) => e.stopPropagation()}
                placeholder="Search the web..."
                className="w-full h-14 px-6 pr-24 text-lg bg-white/70 backdrop-blur-xl border border-gray-200/50 rounded-2xl shadow-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-transparent transition-all duration-300 group-hover:shadow-xl"
              />
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2 flex items-center space-x-2">
                {query && (
                  <button
                    onClick={clearSearch}
                    className="p-1.5 hover:bg-gray-100 rounded-full transition-colors duration-200"
                  >
                    <X className="w-4 h-4 text-gray-400" />
                  </button>
                )}
                <button
                  onClick={() => handleSearch()}
                  disabled={isSearching}
                  className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-xl transition-all duration-200 disabled:opacity-50"
                >
                  {isSearching ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <Search className="w-4 h-4" />
                  )}
                </button>
              </div>
            </div>

            {/* Suggestions Dropdown */}
            {showSuggestions && !searchResults.length && (
              <div className="absolute top-full left-0 right-0 mt-2 bg-white/90 backdrop-blur-xl border border-gray-200/50 rounded-2xl shadow-xl z-40 overflow-hidden">
                {recentSearches.length > 0 && (
                  <div className="p-4 border-b border-gray-100">
                    <div className="flex items-center space-x-2 mb-3">
                      <Clock className="w-4 h-4 text-gray-400" />
                      <span className="text-sm font-medium text-gray-600">Recent</span>
                    </div>
                    {recentSearches.map((search) => (
                      <button
                        key={search.id}
                        onClick={() => handleSuggestionClick(search.query)}
                        className="block w-full text-left px-3 py-2 text-gray-700 hover:bg-gray-50 rounded-lg transition-colors duration-150"
                      >
                        {search.query}
                      </button>
                    ))}
                  </div>
                )}
                
                <div className="p-4">
                  <div className="flex items-center space-x-2 mb-3">
                    <TrendingUp className="w-4 h-4 text-gray-400" />
                    <span className="text-sm font-medium text-gray-600">Trending</span>
                  </div>
                  {trendingTopics.map((topic) => (
                    <button
                      key={topic.id}
                      onClick={() => handleSuggestionClick(topic.title)}
                      className="block w-full text-left px-3 py-2 text-gray-700 hover:bg-gray-50 rounded-lg transition-colors duration-150"
                    >
                      <div className="flex items-center justify-between">
                        <span>{topic.title}</span>
                        <span className="text-xs text-gray-400">
                          {(topic.searches / 1000000).toFixed(1)}M
                        </span>
                      </div>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Search Results */}
        {searchResults.length > 0 && (
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-gray-900">
                Search Results
              </h2>
              <span className="text-sm text-gray-500">
                About {searchResults.length} results
              </span>
            </div>
            
            <div className="space-y-4">
              {searchResults.map((result, index) => (
                <div
                  key={result.id}
                  className="group p-6 bg-white/70 backdrop-blur-xl border border-gray-200/50 rounded-2xl hover:shadow-lg transition-all duration-300 animate-fadeIn"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <div className="flex items-start space-x-4">
                    <div className="flex-shrink-0 w-8 h-8 bg-gradient-to-br from-green-400 to-blue-500 rounded-lg flex items-center justify-center">
                      <div className="w-3 h-3 bg-white rounded-full" />
                    </div>
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-blue-600 group-hover:text-blue-700 transition-colors duration-200 mb-1">
                        {result.title}
                      </h3>
                      <p className="text-sm text-green-600 mb-2">{result.url}</p>
                      <p className="text-gray-600 leading-relaxed">{result.description}</p>
                      <div className="mt-3 flex items-center space-x-4">
                        <span className="text-xs px-2 py-1 bg-gray-100 text-gray-600 rounded-full capitalize">
                          {result.type}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {!searchResults.length && !isSearching && query === '' && (
          <div className="text-center py-16">
            <div className="w-16 h-16 bg-gradient-to-br from-blue-400 to-purple-500 rounded-2xl flex items-center justify-center mx-auto mb-6">
              <Search className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-2xl font-semibold text-gray-900 mb-3">
              Search the web
            </h2>
            <p className="text-gray-600 max-w-md mx-auto">
              Enter your search query above to find relevant information from across the internet.
            </p>
          </div>
        )}
      </main>

      <style>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        
        .animate-fadeIn {
          animation: fadeIn 0.6s ease-out forwards;
        }
      `}</style>
    </div>
  );
};

export default SearchComponent;