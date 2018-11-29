require 'rake'

def path(directory, file)
  file File.join(directory, file) do |t|
    cd directory do
      yield t
    end
  end
end

def default(**kwargs)
  task(**kwargs)
  task default: kwargs.first[0]
end

task :clean do
  sh 'git clean -dfx'
end
