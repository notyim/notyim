module Trinity
  class Config
    class << self
      def configure
        yield instance
      end

      def instance
      end
    end
  end
end